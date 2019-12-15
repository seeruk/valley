package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Equals ...
func Equals(_ interface{}) valley.Constraint {
	return valley.Constraint{}
}

// equalsGenerator ...
func equalsGenerator(ctx valley.Context, _ ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	if len(opts) != 1 {
		return output, errors.New("expected exactly one option")
	}

	value, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	predicate := fmt.Sprintf("%s != %s", ctx.VarName, value)
	message := "values must be equal"
	details := map[string]interface{}{
		"equal_to": fmt.Sprintf("%v", value),
	}

	output.Imports = CollectExprImports(ctx, opts[0])
	output.Code = GenerateStandardConstraint(ctx, predicate, message, details)

	return output, nil
}
