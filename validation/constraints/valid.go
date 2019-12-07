package constraints

import (
	"encoding/json"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Valid ...
func Valid() valley.Constraint {
	return valley.Constraint{}
}

// valid ...
func valid(ctx valley.Context, fieldType ast.Expr, _ json.RawMessage) (valley.ConstraintOutput, error) {
	return valley.ConstraintOutput{}, nil
}
