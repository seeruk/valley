package constraints

import (
	"errors"
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
				"%s": %s,
			},
		}
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
		var predicate string

		if len(opts) != 1 {
			return output, errors.New("expected exactly one option")
		}

		variable := valley.Variable{
			Name:  GenerateVariableName(ctx),
			Value: "WIP",
		}

		_ = variable
		_ = predicate

		return output, nil
	}
}
