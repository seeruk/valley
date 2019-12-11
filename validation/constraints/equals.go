package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Equals ...
func Equals(_ interface{}) valley.Constraint {
	return valley.Constraint{}
}

const equalsFormat = `
	if %[1]s != %[2]s {
		%[3]s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "values must be equal",
			Details: map[string]interface{}{
				"equal_to": %[2]s,
			},
		})
		%[4]s
	}
`

// equals ...
func equals(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	value, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	output.Code = fmt.Sprintf(equalsFormat,
		ctx.VarName,
		value,
		ctx.BeforeViolation,
		ctx.AfterViolation,
	)

	return output, nil
}
