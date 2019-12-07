package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

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

// Required ...
func Required() valley.Constraint {
	return valley.Constraint{}
}

// required ...
func required(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintOutput, error) {
	var output valley.ConstraintOutput

	predicate, err := GenerateEmptinessPredicate(ctx.VarName, fieldType)
	if err != nil {
		// TODO: Wrap.
		return output, err
	}

	output.Code = fmt.Sprintf(requiredFormat,
		predicate,
		ctx.BeforeViolation,
		ctx.AfterViolation,
	)

	return output, nil
}
