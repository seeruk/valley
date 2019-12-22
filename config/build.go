package config

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"github.com/seeruk/valley"
)

// importPath is the import path that is used to import Valley types.
const importPath = "github.com/seeruk/valley"

// BuildFromSource builds Config for all types in a given Source by picking out each type that has a
// constraints method defined (in the same file), and using the body of those methods to produce the
// configuration.
func BuildFromSource(src valley.Source) (valley.Config, error) {
	config := valley.Config{
		Types: make(map[string]valley.TypeConfig),
	}

	constraintsMethods := collectConstraintsMethods(src)

	for typeName, method := range constraintsMethods {
		typeConfig, err := buildTypeConfig(src, method)
		if err != nil {
			return config, err
		}

		config.Types[typeName] = typeConfig
	}

	return config, nil
}

// buildTypeConfig builds TypeConfig based on the body of a constraints method in the given Source.
// It does this by reading the Go AST for the file, and picking out calls that match the expected
// usage for Valley.
func buildTypeConfig(src valley.Source, method valley.Method) (valley.TypeConfig, error) {
	config := valley.TypeConfig{
		Fields: make(map[string]valley.FieldConfig),
	}

	for _, stmt := range method.Body.List {
		chain, ok := buildCallChain(src, method, stmt)
		if !ok {
			continue
		}

		var predicate ast.Expr

		typeMethod := chain.Next

		for typeMethod != nil {
			// Assumed to always succeed at this point:
			typeMethodCall := typeMethod.Call
			typeMethodFunc := typeMethodCall.Fun.(*ast.SelectorExpr)

			// Handle the different possible methods chained off of the Type type.
			switch typeMethodFunc.Sel.Name {
			case "Constraints":
				constraints, err := buildConstraintsCall(src, predicate, typeMethod)
				if err != nil {
					return config, err
				}

				// Merge the result, as multiple separate calls to Constraints could be made.
				config.Constraints = append(config.Constraints, constraints...)

				// Set up the next call, if there is one.
				typeMethod = typeMethod.Next
			case "Field":
				fieldName, fieldConfig, err := buildFieldsCall(src, method, predicate, typeMethod)
				if err != nil {
					return config, err
				}

				// Merge the new configuration with any existing configuration.
				existingConfig := config.Fields[fieldName]
				existingConfig.Constraints = append(existingConfig.Constraints, fieldConfig.Constraints...)
				existingConfig.Elements = append(existingConfig.Elements, fieldConfig.Elements...)
				existingConfig.Keys = append(existingConfig.Keys, fieldConfig.Keys...)
				config.Fields[fieldName] = existingConfig

				// Field doesn't return Type, so there can be no further method calls.
				typeMethod = nil
			case "When":
				if len(typeMethodCall.Args) != 1 {
					return config, errorOn(src, typeMethod.Call.Pos(), "exactly one argument should be passed to When")
				}

				predicate = typeMethodCall.Args[0]
				typeMethod = typeMethod.Next
			}
		}
	}

	return config, nil
}

// buildCallChain ...
func buildCallChain(src valley.Source, method valley.Method, stmt ast.Stmt) (*callExprNode, bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		warnOn(src, stmt.Pos(), "skipping line that is not a statement")
		return nil, false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		// This also protects us later, chain.Next should never be nil after this check.
		warnOn(src, stmt.Pos(), "skipping line that is not a call expression")
		return nil, false
	}

	// Each call expression in a Go AST function body is right-to-left, so the last method call
	// is the first thing you see. We want the opposite, because it's easier to verify what the
	// methods are being called on (i.e. that it's a valley.Type), and then which method is
	// being called to determine how to behave from that point on.
	chain, err := buildCallExpr(src, callExpr)
	if err != nil {
		warnOn(src, stmt.Pos(), "skipping line with unexpected structure: %v", err)
		return nil, false
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
		warnOn(src, stmt.Pos(), "skipping call that isn't on valley.Type")
		return nil, false
	}

	return chain, true
}

// buildConstraintsCall ...
func buildConstraintsCall(src valley.Source, predicate ast.Expr, typeMethod *callExprNode) ([]valley.ConstraintConfig, error) {
	var configs []valley.ConstraintConfig

	for _, expr := range typeMethod.Call.Args {
		constraintConfig, err := buildConstraintConfig(src, predicate, expr)
		if err != nil {
			return nil, err
		}

		configs = append(configs, constraintConfig)
	}

	return configs, nil
}

// buildFieldsCall ...
func buildFieldsCall(src valley.Source, method valley.Method, predicate ast.Expr, typeMethod *callExprNode) (string, valley.FieldConfig, error) {
	var config valley.FieldConfig

	if len(typeMethod.Call.Args) != 1 {
		// No field was passed to Field, one must be given to know which field any
		// constraints apply to; or nothing was chained off of the call to Field, so just
		// continue... nothing to configure at this point.
		return "", config, errorOn(src, typeMethod.Call.Pos(), "exactly one argument should be passed to Field")
	}

	if typeMethod.Next == nil {
		return "", config, errorOn(src, typeMethod.Call.Pos(), "a method should be called on Field")
	}

	fieldArg, ok := typeMethod.Call.Args[0].(*ast.SelectorExpr)
	if !ok {
		// The argument passed to Field must be a selector (i.e. a field on the type).
		return "", config, errorOn(src, typeMethod.Call.Pos(), "value passed to Field should be a selector")
	}

	fieldArgOn, ok := fieldArg.X.(*ast.Ident)
	if !ok || fieldArgOn.Name != method.Receiver {
		// The argument passed to Field must be on an ident (i.e. the type). Additionally,
		// the argument passed to Field must be on the receiver for the constraints method.
		return "", config, errorOn(src, fieldArg.Pos(), "value passed to Field should be a field on the receiver's type")
	}

	fieldConfig, err := buildFieldConfig(src, predicate, typeMethod.Next)
	if err != nil {
		return "", config, err
	}

	return fieldArg.Sel.Name, fieldConfig, nil
}

// buildFieldConfig ...
func buildFieldConfig(src valley.Source, predicate ast.Expr, fieldMethodNode *callExprNode) (valley.FieldConfig, error) {
	var config valley.FieldConfig

	for _, expr := range fieldMethodNode.Call.Args {
		// The "expr" here is the argument being passed to a method on the valley.Field type, in
		// other words, we expect each of these arguments to be a constraint.
		constraintConfig, err := buildConstraintConfig(src, predicate, expr)
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
		case "Keys":
			config.Keys = append(config.Keys, constraintConfig)
		}
	}

	if fieldMethodNode.Next != nil {
		nextConfig, err := buildFieldConfig(src, predicate, fieldMethodNode.Next)
		if err != nil {
			return config, err
		}

		config.Constraints = append(config.Constraints, nextConfig.Constraints...)
		config.Elements = append(config.Elements, nextConfig.Elements...)
		config.Keys = append(config.Keys, nextConfig.Keys...)
	}

	return config, nil
}

// buildCallExpr converts the chain of Go AST expressions for a field call statement into a linked
// list of each node. This can later be reversed to get the calls in left-to-right order which is
// easier to validate (i.e. check if the call is on the valley.Type).
func buildCallExpr(src valley.Source, outer ast.Expr) (*callExprNode, error) {
	if ident, ok := outer.(*ast.Ident); ok {
		return &callExprNode{
			Ident: ident,
		}, nil
	}

	callExpr, ok := outer.(*ast.CallExpr)
	if !ok {
		return nil, errorOn(src, outer.Pos(), "statement expression must be a call")
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, errorOn(src, outer.Pos(), "statement expression must be a method call")
	}

	next, err := buildCallExpr(src, selectorExpr.X)
	if err != nil {
		return nil, err
	}

	return &callExprNode{
		Call: callExpr,
		Next: next,
	}, nil
}

// buildConstraintConfig ...
func buildConstraintConfig(src valley.Source, predicate, expr ast.Expr) (valley.ConstraintConfig, error) {
	var config valley.ConstraintConfig

	config.Predicate = predicate

	constraintCall, ok := expr.(*ast.CallExpr)
	if !ok {
		// If the argument to the Constraints call wasn't a function call, it's invalid.
		return config, errorOn(src, expr.Pos(), "constraint must be a function call")
	}

	constraintFunc, ok := constraintCall.Fun.(*ast.SelectorExpr)
	if !ok {
		// If the function call wasn't on a selector (i.e. a function in a package is what we're
		// looking for), then maybe it's a local constraint. We don't support that yet.
		// TODO: Is this limitation necessary? It might complicate this code a bit.
		return config, errorOn(src, expr.Pos(), "constraint must be from a different package")
	}

	constraintFuncOn, ok := constraintFunc.X.(*ast.Ident)
	if !ok {
		// If the function call was a selector, but wasn't on an ident (i.e. we're assuming it
		// should be a package name / alias).
		return config, errorOn(src, expr.Pos(), "constraints must be exported functions in a package")
	}

	constraintFuncPkg, ok := findImportByName(src.Imports, constraintFuncOn.Name)
	if !ok {
		// If the function call was on an ident, but it wasn't a package that was imported (or if
		// the package name is different to the alias, or the end of the import path).
		return config, errorOn(src, expr.Pos(), "constraint must be defined in an imported package")
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

// collectConstraintsMethods looks through the methods in a given Source and extracts the first
// method that looks like a constraints method. This allows the Constraint method to have any name.
//
// TODO: Having any name is not very useful, it's only useful if you can define more than one, and
// have each constraints method generate a different validation function at the end.
func collectConstraintsMethods(src valley.Source) map[string]valley.Method {
	constraintsMethods := make(map[string]valley.Method)

	for typeName, methods := range src.Methods {
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

			imp, ok := findImportByName(src.Imports, selectorPkg.Name)
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
func findImportByName(imports []valley.Import, name string) (valley.Import, bool) {
	for _, imp := range imports {
		if imp.Alias == name {
			return imp, true
		}
	}

	return valley.Import{}, false
}

// errorOn returns an error with the given message, in the given Source, at the given position.
func errorOn(src valley.Source, pos token.Pos, message string, args ...interface{}) error {
	return errors.New(messageOn(src, pos, message, args...))
}

// warnOn prints a given warning message, in the given Source, at the given position.
func warnOn(src valley.Source, pos token.Pos, message string, args ...interface{}) {
	fmt.Printf("valley: %s\n", messageOn(src, pos, message, args...))
}

// messageOn returns a formatted message, in the given Source, at the given position.
func messageOn(src valley.Source, pos token.Pos, message string, args ...interface{}) string {
	position := src.FileSet.Position(pos)

	srcPath, err := filepath.Abs(src.FileName)
	if err != nil {
		srcPath = src.FileName
	}

	args = append(args, position.Line, position.Column, srcPath)

	return fmt.Sprintf(message+" on line %d, col %d in '%s'", args...)
}
