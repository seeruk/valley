package valley

import (
	"fmt"
	"go/ast"
	"go/token"
)

// importPath is the import path that is used to import Valley types.
const importPath = "github.com/seeruk/valley/valley"

// Config ...
type Config struct {
	Types map[string]TypeConfig `json:"types"`
}

// TypeConfig ...
type TypeConfig struct {
	Constraints []ConstraintConfig     `json:"constraints"`
	Fields      map[string]FieldConfig `json:"fields"`
}

// FieldConfig ...
type FieldConfig struct {
	Constraints []ConstraintConfig `json:"constraints"`
	Elements    []ConstraintConfig `json:"elements"`
}

// ConstraintConfig ...
type ConstraintConfig struct {
	Name string     `json:"name"`
	Opts []ast.Expr `json:"opts"`
	Pos  token.Pos
}

// ConfigFromFile builds Config for all types in a given file by picking out each type that has
// a constraints method defined (in the same file), and using the body of those methods to produce
// the configuration.
func ConfigFromFile(file File) (Config, error) {
	config := Config{
		Types: make(map[string]TypeConfig),
	}

	constraintsMethods := collectConstraintsMethod(file)

	for typeName, method := range constraintsMethods {
		config.Types[typeName] = buildTypeConfig(file, method)
	}

	return config, nil
}

// buildTypeConfig builds TypeConfig based on the body of a constraints method in the given file.
// It does this by reading the Go AST for the file, and picking out calls that match the expected
// usage for Valley.
func buildTypeConfig(file File, method Method) TypeConfig {
	config := TypeConfig{
		Fields: make(map[string]FieldConfig),
	}

	for _, stmt := range method.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}

		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		// Each call expression in a Go AST function body is right-to-left, so the last method call
		// is the first thing you see. We want the opposite, because it's easier to verify what the
		// methods are being called on (i.e. that it's a valley.Type), and then which method is
		// being called to determine how to behave from that point on.
		chain := buildCallExpr(callExpr).Reverse()

		// At this point we should know there is only one parameter, and it should be the
		// valley.Type argument.
		param := method.Params.List[0]
		paramName := param.Names[0].Name

		if chain.Ident == nil || chain.Ident.Name != paramName {
			// The call should be happening on the valley.Type. It doesn't have to be called `t`, so
			// we get the parameter name and compare it to the root identifier of the call chain
			// (i.e. the identifier all of the calls in the statement are coming off of).
			errorOn(file, stmt.Pos(), "skipping call that isn't on valley.Type")
			continue
		}

		if chain.Next == nil || chain.Next.Call == nil {
			// No call after ident, i.e. a lone reference to the valley.Type? Can this happen?
			continue
		}

		typeMethod := chain.Next
		typeMethodCall := typeMethod.Call
		typeMethodFunc, ok := typeMethodCall.Fun.(*ast.SelectorExpr)
		if !ok {
			// This should probably never happen given we know we're calling a method on valley.Type
			// at this point. It should always have a selector, and it should always be the Type.
			continue
		}

		// TODO: How much stuff can we trust above to remove some code above?
		// TODO: This next bit should go elsewhere, and Constraints should return the valley.Type so
		// we can chain calls, and still handle that.

		switch typeMethodFunc.Sel.Name {
		case "Constraints":
			for _, expr := range chain.Next.Call.Args {
				constraintCall, ok := expr.(*ast.CallExpr)
				if !ok {
					// If the argument to the Constraints call wasn't a function call, it's invalid.
					continue
				}

				constraintFunc, ok := constraintCall.Fun.(*ast.SelectorExpr)
				if !ok {
					// If the function call wasn't on a selector (i.e. a function in a package is
					// what we're looking for). This could still pass if we're calling a method
					// returned by a function that returns another type, so continue...
					continue
				}

				constraintFuncOn, ok := constraintFunc.X.(*ast.Ident)
				if !ok {
					// If the function call was a selector, but wasn't on an ident (i.e. we're
					// assuming it should be a package name / alias.
					continue
				}

				constraintFuncPkg, ok := findImportByName(file.Imports, constraintFuncOn.Name)
				if !ok {
					// If the function call was on an ident, but it wasn't a package that was
					// imported (or if the package name is different to the alias, or the end of the
					// import path).
					continue
				}

				config.Constraints = append(config.Constraints, ConstraintConfig{
					Name: fmt.Sprintf("%s.%s", constraintFuncPkg.Path, constraintFunc.Sel.Name),
					Opts: constraintCall.Args,
					Pos:  expr.Pos(),
				})
			}
		case "Field":
			if len(typeMethodCall.Args) != 1 || typeMethod.Next == nil {
				// No field was passed to Field, one must be given to know which field any
				// constraints apply to; or nothing was chained off of the call to Field, so just
				// continue... nothing to configure at this point.
				continue
			}

			fieldArg, ok := typeMethodCall.Args[0].(*ast.SelectorExpr)
			if !ok {
				// The argument passed to Field must be a selector (i.e. a field on the type).
				continue
			}

			fieldArgOn, ok := fieldArg.X.(*ast.Ident)
			if !ok || fieldArgOn.Name != method.Receiver {
				// The argument passed to Field must be on an ident (i.e. the type). Additionally,
				// the argument passed to Field must be on the receiver for the constraints method.
				continue
			}

			config.Fields[fieldArg.Sel.Name] = buildFieldConfig(file, typeMethod.Next)
		}
	}

	return config
}

// buildFieldConfig ...
func buildFieldConfig(file File, fieldMethodNode *callExprNode) FieldConfig {
	var config FieldConfig

	for _, argExpr := range fieldMethodNode.Call.Args {
		// The "arg" here is the argument being passed to a method on the valley.Field type, in
		// other words, we expect each of these arguments to be a constraint, so we verify that.

		argCall, ok := argExpr.(*ast.CallExpr)
		if !ok {
			// Argument was not a function call, all constraints are, so this is not a constraint.
			continue
		}

		argFunc, ok := argCall.Fun.(*ast.SelectorExpr)
		if !ok {
			// Argument was a function call, but the function call was not a selector (e.g. this
			// will be false if the function call was from the same package, a local function).
			// Currently because of the way the constraints are registered, they must all be in a
			// separate package.
			continue
		}

		// This assumes everything is a function in a package?
		argFuncOn, ok := argFunc.X.(*ast.Ident)
		if !ok {
			// If the function was chained off of anything other than an ident (e.g. the result of
			// another function), then it's not currently supported.
			continue
		}

		argFuncPkg, ok := findImportByName(file.Imports, argFuncOn.Name)
		if !ok {
			// If the function was chained off of anything other than an import that exists in the
			// file (where the package name matches the end of the import path, or the import
			// alias), for example a package-local variable, then it's currently not supported.
			continue
		}

		// This should be one of the methods on the valley.Field type.
		fieldMethodFunc, ok := fieldMethodNode.Call.Fun.(*ast.SelectorExpr)
		if !ok {
			// TODO: Node itself is not a function call (i.e. the thing args are being passed into.
			continue
		}

		constraintConfig := ConstraintConfig{
			Name: fmt.Sprintf("%s.%s", argFuncPkg.Path, argFunc.Sel.Name),
			Opts: argCall.Args,
			Pos:  argExpr.Pos(),
		}

		switch fieldMethodFunc.Sel.Name {
		case "Constraints":
			config.Constraints = append(config.Constraints, constraintConfig)
		case "Elements":
			config.Elements = append(config.Elements, constraintConfig)
		}
	}

	if fieldMethodNode.Next != nil {
		nextConfig := buildFieldConfig(file, fieldMethodNode.Next)

		config.Constraints = append(config.Constraints, nextConfig.Constraints...)
		config.Elements = append(config.Elements, nextConfig.Elements...)
	}

	return config
}

// buildCallExpr converts the chain of Go AST expressions for a field call statement into a linked
// list of each node. This can later be reversed to get the calls in left-to-right order which is
// easier to validate (i.e. check if the call is on the valley.Type).
func buildCallExpr(outer ast.Expr) *callExprNode {
	node := &callExprNode{}

	if ident, ok := outer.(*ast.Ident); ok {
		node.Ident = ident
	}

	if callExpr, ok := outer.(*ast.CallExpr); ok {
		node.Call = callExpr
	}

	callExpr, ok := outer.(*ast.CallExpr)
	if !ok {
		return node
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return node
	}

	node.Next = buildCallExpr(selectorExpr.X)

	return node
}

// callExprNode is a linked list node for the components that make up a Go AST statement expression
// of a method call (potentially chained, i.e. method call, from method call, from method call, on a
// variable).
type callExprNode struct {
	Ident *ast.Ident
	Call  *ast.CallExpr
	Next  *callExprNode
}

// Reverse flips the order the nodes from this node in the list.
func (n *callExprNode) Reverse() *callExprNode {
	current := n

	var prev *callExprNode
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// collectConstraintsMethod looks through the methods in a given file and extracts the first method
// that looks like a constraints method. This allows the Constraint method to have any name.
//
// TODO: Having any name is not very useful, it's only useful if you can define more than one, and
// have each constraints method generate a different validation function at the end.
func collectConstraintsMethod(file File) map[string]Method {
	constraintsMethods := make(map[string]Method)

	for typeName, methods := range file.Methods {
		for _, method := range methods {
			if method.Results != nil || method.Params == nil || len(method.Params.List) != 1 {
				// Valley constraints methods don't return anything, and have one param.
				continue
			}

			param := method.Params.List[0]

			selector, ok := param.Type.(*ast.SelectorExpr)
			if !ok || selector.Sel.Name != "Type" {
				// The type name must be what we expect.
				continue
			}

			selectorPkg, ok := selector.X.(*ast.Ident)
			if !ok {
				continue
			}

			imp, ok := findImportByName(file.Imports, selectorPkg.Name)
			if !ok || imp.Path != importPath {
				// The type must come from our code!
				continue
			}

			constraintsMethods[typeName] = method
			break
		}
	}

	return constraintsMethods
}

// findImportByName looks for an import with the given name (or alias) in the given set of imports.
func findImportByName(imports []Import, alias string) (Import, bool) {
	for _, imp := range imports {
		if imp.Alias == alias {
			return imp, true
		}
	}

	return Import{}, false
}

// errorOn prints a given error message, in the given file, at the given position.
func errorOn(file File, pos token.Pos, message string, args ...interface{}) {
	position := file.FileSet.Position(pos)

	args = append(args, position.Line, position.Column, file.Package, file.Name)

	// TODO: I have no doubt this could be more robust and useful. Maybe this could the filename
	// from the root of the currently module path instead? Would need to get CWD, module path,
	// figure out where we are, and put it all together to get the path from the root of the module
	// to the file.
	fmt.Printf("valley: "+message+" on line %d, col %d in %s/%s\n", args...)
}
