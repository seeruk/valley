package constraints

import (
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Valid ...
func Valid() valley.Constraint {
	return valley.Constraint{}
}

// valid ...
func valid(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintOutput, error) {
	return valley.ConstraintOutput{}, nil
}
