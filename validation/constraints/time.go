package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"time"

	"github.com/seeruk/valley"
)

// TimeAfter ...
func TimeAfter(after time.Time) valley.Constraint {
	return valley.Constraint{}
}

// TimeBefore ...
func TimeBefore(before time.Time) valley.Constraint {
	return valley.Constraint{}
}

// timeFormat is the format used for rendering a `TimeAfter` or `TimeBefore` constraint.
const timeFormat = `
	if %s {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "%s",
			Details: map[string]interface{}{
				"time": %s,
			},
		})
		%s
	}
`

// Possible timeKind values.
const (
	after  timeKind = "after"
	before timeKind = "before"
)

// timeKind ...
type timeKind string

// timeGenerator ...
func timeGenerator(kind timeKind) valley.ConstraintGenerator {
	return func(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
		var output valley.ConstraintGeneratorOutput
		var message string
		var predicate string

		if len(opts) != 1 {
			return output, errors.New("expected exactly one option")
		}

		output.Imports = CollectExprImports(ctx, opts[0])

		timeSelector, err := SprintNode(ctx.Source.FileSet, opts[0])
		if err != nil {
			return output, fmt.Errorf("failed to render expression: %v", err)
		}

		_, isPointer := fieldType.(*ast.StarExpr)

		varName := ctx.VarName
		if isPointer {
			predicate += fmt.Sprintf("%s != nil && ", varName)
			varName = "*" + varName
		}

		// TODO: These messages aren't great - any way to improve them?
		if kind == after {
			message = "value must be after time"
			predicate += fmt.Sprintf("!%s.After(%s)", varName, timeSelector)
		} else {
			message = "value must be before time"
			predicate += fmt.Sprintf("!%s.Before(%s)", varName, timeSelector)
		}

		output.Code = fmt.Sprintf(timeFormat,
			predicate,
			ctx.BeforeViolation,
			message,
			timeSelector,
			ctx.AfterViolation,
		)

		return output, timeTypeCheck(fieldType)
	}
}

// timeTypeCheck ...
func timeTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return timeTypeCheck(e)
	case *ast.SelectorExpr:
		return nil
	}

	return ErrTypeWarning
}
