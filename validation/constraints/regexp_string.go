package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Regexp ...
func RegexpString(regexp string) valley.Constraint {
	return valley.Constraint{}
}

// regexpStringGenerator ...
func regexpStringGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var predicate string

	if len(opts) != 1 {
		return output, errors.New("expected exactly one option")
	}

	pattern, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	patternVarName := GenerateVariableName(ctx)

	// Add the regexp as a variable, that will be compiled when imported - should help performance.
	output.Vars = []valley.Variable{
		{Name: patternVarName, Value: fmt.Sprintf("regexp.MustCompile(%s)", pattern)},
	}

	// Check if the field is a pointer, if so, we'll add a nil check and dereference from there.
	_, isPointer := fieldType.(*ast.StarExpr)

	varName := ctx.VarName
	if isPointer {
		predicate += fmt.Sprintf("%s != nil && ", varName)
		varName = "*" + varName
	}

	predicate += fmt.Sprintf("!%s.MatchString(%s)", patternVarName, varName)
	message := "value must match regular expression"
	details := map[string]interface{}{
		"regexp": fmt.Sprintf("%s.String()", patternVarName),
	}

	output.Imports = CollectExprImports(ctx, opts[0])
	output.Imports = append(output.Imports, valley.Import{
		Path:  "regexp",
		Alias: "regexp",
	})

	output.Code = GenerateStandardConstraint(ctx, predicate, message, details)

	return output, regexpStringTypeCheck(fieldType)
}

// regexpStringTypeCheck ...
func regexpStringTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return regexpStringTypeCheck(e.X)
	case *ast.Ident:
		if e.Name == "string" {
			return nil
		}
	}

	return ErrTypeWarning
}
