package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Required ...
func Required() valley.Constraint {
	return valley.Constraint{}
}

// requiredFormat is the format used for rendering a `required` constraint.
const requiredFormat = `
	if %s {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})
		%s
	}
`

// required ...
func required(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintOutput, error) {
	var output valley.ConstraintOutput

	predicate, err := GenerateEmptinessPredicate(ctx.VarName, fieldType)
	if err != nil {
		return output, err
	}

	output.Code = fmt.Sprintf(requiredFormat,
		predicate,
		ctx.BeforeViolation,
		ctx.AfterViolation,
	)

	return output, nil
}
