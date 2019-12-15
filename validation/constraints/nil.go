package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Nil ...
func Nil() valley.Constraint {
	return valley.Constraint{}
}

// nilGenerator ...
func nilGenerator(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	return valley.ConstraintGeneratorOutput{
		Code: GenerateStandardConstraint(ctx,
			fmt.Sprintf("%s != nil", ctx.VarName),
			"value must be nil",
			nil,
		),
	}, nilTypeCheck(fieldType)
}

// nilTypeCheck ...
func nilTypeCheck(expr ast.Expr) error {
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
