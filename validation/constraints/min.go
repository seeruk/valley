package constraints

import (
	"encoding/json"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Min ...
func Min(min float64) valley.Constraint {
	return valley.Constraint{}
}

// min ...
func min(ctx valley.Context, fieldType ast.Expr, _ json.RawMessage) (valley.ConstraintOutput, error) {
	return valley.ConstraintOutput{}, nil
}
