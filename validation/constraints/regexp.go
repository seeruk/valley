package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"regexp"

	"github.com/seeruk/valley"
)

// Regexp ...
func Regexp(regexp *regexp.Regexp) valley.Constraint {
	return valley.Constraint{}
}

// regexpStringGenerator ...
func regexpGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var predicate string

	if len(opts) != 1 {
		return output, errors.New("expected exactly one option")
	}

	patternSelector, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	// Check if the field is a pointer, if so, we'll add a nil check and dereference from there.
	_, isPointer := fieldType.(*ast.StarExpr)

	varName := ctx.VarName
	if isPointer {
		predicate += fmt.Sprintf("%s != nil && ", varName)
		varName = "*" + varName
	}

	predicate += fmt.Sprintf("!%s.MatchString(%s)", patternSelector, varName)
	message := "value must match regular expression"
	details := map[string]interface{}{
		"regexp": fmt.Sprintf("%s.String()", patternSelector),
	}

	output.Imports = CollectExprImports(ctx, opts[0])
	output.Code = GenerateStandardConstraint(ctx, predicate, message, details)

	return output, regexpTypeCheck(fieldType)
}

// regexpStringTypeCheck ...
func regexpTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return regexpTypeCheck(e.X)
	case *ast.SelectorExpr:
		return nil
	case *ast.Ident:
		if e.Name == "string" {
			return nil
		}
	}

	return ErrTypeWarning
}
