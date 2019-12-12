package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// NotNil ...
func NotNil() valley.Constraint {
	return valley.Constraint{}
}

const notNilFormat = `
	if %s == nil {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "value must not be nil",
		})
		%s
	}
`

// notNil ...
func notNil(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	return valley.ConstraintGeneratorOutput{
		Code: fmt.Sprintf(notNilFormat,
			ctx.VarName,
			ctx.BeforeViolation,
			ctx.AfterViolation,
		),
	}, notNilTypeCheck(fieldType)
}

// notNilTypeCheck ...
func notNilTypeCheck(expr ast.Expr) error {
	switch expr.(type) {
	case *ast.StarExpr:
		return nil
	case *ast.ArrayType:
		return nil
	case *ast.ChanType:
		return nil
	case *ast.MapType:
		return nil
	}

	return ErrTypeWarning
}
