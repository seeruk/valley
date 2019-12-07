package constraints

import (
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Min ...
func Min(min float64) valley.Constraint {
	return valley.Constraint{}
}

// min ...
func min(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintOutput, error) {
	return valley.ConstraintOutput{}, nil
}
