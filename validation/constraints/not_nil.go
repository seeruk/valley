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

// notNilGenerator ...
func notNilGenerator(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	return valley.ConstraintGeneratorOutput{
		Code: GenerateStandardConstraint(ctx,
			fmt.Sprintf("%s == nil", ctx.VarName),
			"value must not be nil",
			nil,
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
