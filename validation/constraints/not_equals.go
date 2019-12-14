package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// NotEquals ...
func NotEquals(_ interface{}) valley.Constraint {
	return valley.Constraint{}
}

const notEqualsFormat = `
	if %[1]s == %[2]s {
		%[3]s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "values must not be equal",
			Details: map[string]interface{}{
				"equal_to": %[2]s,
			},
		})
		%[4]s
	}
`

// notEqualsGenerator ...
func notEqualsGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	if len(opts) != 1 {
		return output, errors.New("expected exactly one option")
	}

	value, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	output.Imports = CollectExprImports(ctx, opts[0])

	output.Code = fmt.Sprintf(notEqualsFormat,
		ctx.VarName,
		value,
		ctx.BeforeViolation,
		ctx.AfterViolation,
	)

	return output, nil
}
