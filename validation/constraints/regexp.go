package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	regex "regexp"

	"github.com/seeruk/valley"
)

// Regexp ...
func Regexp(regexp *regex.Regexp) valley.Constraint {
	return valley.Constraint{}
}

const patternFormat = `
	if %s {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "value must match regular expression",
			Details: map[string]interface{}{
				"regexp": %s.String(),
			},
		})
		%s
	}
`

// regexp ...
func regexp(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var predicate string

	if len(opts) != 1 {
		return output, errors.New("expected exactly one option")
	}

	output.Imports = CollectExprImports(ctx, opts[0])

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

	output.Code = fmt.Sprintf(patternFormat,
		predicate,
		ctx.BeforeViolation,
		patternSelector,
		ctx.AfterViolation,
	)

	return output, regexpTypeCheck(fieldType)
}

// regexpTypeCheck ...
func regexpTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return regexpTypeCheck(e.X)
	case *ast.Ident:
		if e.Name == "string" {
			return nil
		}
	}

	return ErrTypeWarning
}
