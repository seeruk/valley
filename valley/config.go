package valley

import (
	"errors"
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

	constraintsMethods := collectConstraintsMethods(file)

	for typeName, method := range constraintsMethods {
		typeConfig, err := buildTypeConfig(file, method)
		if err != nil {
			// TODO: Wrap?
			return config, err
		}

		config.Types[typeName] = typeConfig
	}

	return config, nil
}

// buildTypeConfig builds TypeConfig based on the body of a constraints method in the given file.
// It does this by reading the Go AST for the file, and picking out calls that match the expected
// usage for Valley.
func buildTypeConfig(file File, method Method) (TypeConfig, error) {
	config := TypeConfig{
		Fields: make(map[string]FieldConfig),
	}

	for _, stmt := range method.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			warnOn(file, stmt.Pos(), "skipping line that is not a statement")
			continue
		}

		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			// This also protects us later, chain.Next should never be nil after this check.
			warnOn(file, stmt.Pos(), "skipping line that is not a call expression")
			continue
		}

		// Each call expression in a Go AST function body is right-to-left, so the last method call
		// is the first thing you see. We want the opposite, because it's easier to verify what the
		// methods are being called on (i.e. that it's a valley.Type), and then which method is
		// being called to determine how to behave from that point on.
		chain, err := buildCallExpr(file, callExpr)
		if err != nil {
			warnOn(file, stmt.Pos(), "skipping line with unexpected structure: %v", err)
			continue
		}

		chain = chain.Reverse()

		// At this point we should know there is only one parameter, and it should be the
		// valley.Type argument.
		param := method.Params.List[0]
		paramName := param.Names[0].Name

		if chain.Ident == nil || chain.Ident.Name != paramName {
			// The call should be happening on the valley.Type. It doesn't have to be called `t`, so
			// we get the parameter name and compare it to the root identifier of the call chain
			// (i.e. the identifier all of the calls in the statement are coming off of).
			warnOn(file, stmt.Pos(), "skipping call that isn't on valley.Type")
			continue
		}

		typeMethod := chain.Next

		for typeMethod != nil {
			typeMethodCall := typeMethod.Call
			typeMethodFunc, ok := typeMethodCall.Fun.(*ast.SelectorExpr)
			if !ok {
				// This should probably never happen given we know we're calling a method on valley.Type
				// at this point. It should always have a selector, and it should always be the Type.
				continue
			}

			// Handle the different possible methods chained off of the Type type.
			switch typeMethodFunc.Sel.Name {
			case "Constraints":
				constraints, err := buildConstraintsCall(file, typeMethod)
				if err != nil {
					return config, err
				}

				// Merge the result, as multiple separate calls to Constraints could be made.
				config.Constraints = append(config.Constraints, constraints...)

				// Set up the next call, if there is one.
				typeMethod = typeMethod.Next
			case "Field":
				fieldName, fieldConfig, err := buildFieldsCall(file, method, typeMethod)
				if err != nil {
					return config, err
				}

				// Merge the new configuration with any existing configuration.
				existingConfig := config.Fields[fieldName]
				existingConfig.Constraints = append(existingConfig.Constraints, fieldConfig.Constraints...)
				existingConfig.Elements = append(existingConfig.Elements, fieldConfig.Elements...)
				config.Fields[fieldName] = existingConfig

				// Field doesn't return Type, so there can be no further method calls.
				typeMethod = nil
			}
		}
	}

	return config, nil
}

// buildConstraintsCall ...
func buildConstraintsCall(file File, typeMethod *callExprNode) ([]ConstraintConfig, error) {
	var configs []ConstraintConfig

	for _, expr := range typeMethod.Call.Args {
		constraintConfig, err := buildConstraintConfig(file, expr)
		if err != nil {
			return nil, err
		}

		configs = append(configs, constraintConfig)
	}

	return configs, nil
}

// buildFieldsCall ...
func buildFieldsCall(file File, method Method, typeMethod *callExprNode) (string, FieldConfig, error) {
	var config FieldConfig

	if len(typeMethod.Call.Args) != 1 {
		// No field was passed to Field, one must be given to know which field any
		// constraints apply to; or nothing was chained off of the call to Field, so just
		// continue... nothing to configure at this point.
		return "", config, errorOn(file, typeMethod.Call.Pos(), "exactly one argument should be passed to Field")
	}

	if typeMethod.Next == nil {
		return "", config, errorOn(file, typeMethod.Call.Pos(), "a method should be called on Field")
	}

	fieldArg, ok := typeMethod.Call.Args[0].(*ast.SelectorExpr)
	if !ok {
		// The argument passed to Field must be a selector (i.e. a field on the type).
		return "", config, errorOn(file, fieldArg.Pos(), "value passed to Field should selector")
	}

	fieldArgOn, ok := fieldArg.X.(*ast.Ident)
	if !ok || fieldArgOn.Name != method.Receiver {
		// The argument passed to Field must be on an ident (i.e. the type). Additionally,
		// the argument passed to Field must be on the receiver for the constraints method.
		return "", config, errorOn(file, fieldArg.Pos(), "value passed to Field should field on the receiver's type")
	}

	fieldConfig, err := buildFieldConfig(file, typeMethod.Next)
	if err != nil {
		return "", config, err
	}

	return fieldArg.Sel.Name, fieldConfig, nil
}

// buildFieldConfig ...
func buildFieldConfig(file File, fieldMethodNode *callExprNode) (FieldConfig, error) {
	var config FieldConfig

	for _, expr := range fieldMethodNode.Call.Args {
		// The "expr" here is the argument being passed to a method on the valley.Field type, in
		// other words, we expect each of these arguments to be a constraint.
		constraintConfig, err := buildConstraintConfig(file, expr)
		if err != nil {
			return config, err
		}

		// This should be one of the methods on the valley.Field type.
		// NOTE: This shouldn't fail, we verify this when building the callExprNode chain.
		fieldMethodFunc, _ := fieldMethodNode.Call.Fun.(*ast.SelectorExpr)

		switch fieldMethodFunc.Sel.Name {
		case "Constraints":
			config.Constraints = append(config.Constraints, constraintConfig)
		case "Elements":
			config.Elements = append(config.Elements, constraintConfig)
		}
	}

	if fieldMethodNode.Next != nil {
		nextConfig, err := buildFieldConfig(file, fieldMethodNode.Next)
		if err != nil {
			return config, err
		}

		config.Constraints = append(config.Constraints, nextConfig.Constraints...)
		config.Elements = append(config.Elements, nextConfig.Elements...)
	}

	return config, nil
}

// buildCallExpr converts the chain of Go AST expressions for a field call statement into a linked
// list of each node. This can later be reversed to get the calls in left-to-right order which is
// easier to validate (i.e. check if the call is on the valley.Type).
func buildCallExpr(file File, outer ast.Expr) (*callExprNode, error) {
	if ident, ok := outer.(*ast.Ident); ok {
		return &callExprNode{
			Ident: ident,
		}, nil
	}

	callExpr, ok := outer.(*ast.CallExpr)
	if !ok {
		return nil, errorOn(file, outer.Pos(), "statement expression must be a call")
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, errorOn(file, outer.Pos(), "statement expression must be a method call")
	}

	next, err := buildCallExpr(file, selectorExpr.X)
	if err != nil {
		// TODO: Wrap?
		return nil, err
	}

	return &callExprNode{
		Call: callExpr,
		Next: next,
	}, nil
}

// buildConstraintConfig ...
func buildConstraintConfig(file File, expr ast.Expr) (ConstraintConfig, error) {
	var config ConstraintConfig

	constraintCall, ok := expr.(*ast.CallExpr)
	if !ok {
		// If the argument to the Constraints call wasn't a function call, it's invalid.
		return config, errorOn(file, expr.Pos(), "constraint must be a function call")
	}

	constraintFunc, ok := constraintCall.Fun.(*ast.SelectorExpr)
	if !ok {
		// If the function call wasn't on a selector (i.e. a function in a package is what we're
		// looking for), then maybe it's a local constraint. We don't support that yet.
		// TODO: Is this limitation necessary? It might complicate this code a bit.
		return config, errorOn(file, expr.Pos(), "constraint must be from a different package")
	}

	constraintFuncOn, ok := constraintFunc.X.(*ast.Ident)
	if !ok {
		// If the function call was a selector, but wasn't on an ident (i.e. we're assuming it
		// should be a package name / alias).
		return config, errorOn(file, expr.Pos(), "constraints must be exported functions in a package")
	}

	constraintFuncPkg, ok := findImportByName(file.Imports, constraintFuncOn.Name)
	if !ok {
		// If the function call was on an ident, but it wasn't a package that was imported (or if
		// the package name is different to the alias, or the end of the import path).
		return config, errorOn(file, expr.Pos(), "constraint must be defined in an imported package")
	}

	config.Name = fmt.Sprintf("%s.%s", constraintFuncPkg.Path, constraintFunc.Sel.Name)
	config.Opts = constraintCall.Args
	config.Pos = expr.Pos()

	return config, nil
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

// collectConstraintsMethods looks through the methods in a given file and extracts the first method
// that looks like a constraints method. This allows the Constraint method to have any name.
//
// TODO: Having any name is not very useful, it's only useful if you can define more than one, and
// have each constraints method generate a different validation function at the end.
func collectConstraintsMethods(file File) map[string]Method {
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
func findImportByName(imports []Import, name string) (Import, bool) {
	for _, imp := range imports {
		if imp.Alias == name {
			return imp, true
		}
	}

	return Import{}, false
}

// errorOn returns an error with the given message, in the given file, at the given position.
func errorOn(file File, pos token.Pos, message string, args ...interface{}) error {
	return errors.New(messageOn(file, pos, message, args...))
}

// warnOn prints a given warning message, in the given file, at the given position.
func warnOn(file File, pos token.Pos, message string, args ...interface{}) {
	fmt.Printf("valley: %s\n", messageOn(file, pos, message, args...))
}

// messageOn returns a formatted message, in the given file, at the given position.
func messageOn(file File, pos token.Pos, message string, args ...interface{}) string {
	position := file.FileSet.Position(pos)

	args = append(args, position.Line, position.Column, file.Package, file.Name)

	// TODO: I have no doubt this could be more robust and useful. Maybe this could the filename
	// from the root of the currently module path instead? Would need to get CWD, module path,
	// figure out where we are, and put it all together to get the path from the root of the module
	// to the file.
	return fmt.Sprintf(message+" on line %d, col %d in '%s/%s'", args...)
}
