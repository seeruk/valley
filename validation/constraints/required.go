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

// requiredFormat is the format used for rendering a `Required` constraint.
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
func required(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	predicate, imports := GenerateEmptinessPredicate(ctx.VarName, fieldType)
	return valley.ConstraintGeneratorOutput{
		Imports: imports,
		Code: fmt.Sprintf(requiredFormat,
			predicate,
			ctx.BeforeViolation,
			ctx.AfterViolation,
		),
	}, nil
}
