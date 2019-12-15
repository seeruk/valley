package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// DeepEquals ...
func DeepEquals(_ interface{}) valley.Constraint {
	return valley.Constraint{}
}

// deepEqualsGenerator ...
func deepEqualsGenerator(ctx valley.Context, _ ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	value, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	predicate := fmt.Sprintf("!reflect.DeepEqual(%s, %s)", ctx.VarName, value)
	message := "values must be deeply equal"
	details := map[string]interface{}{
		"deeply_equal_to": fmt.Sprintf("%v", value),
	}

	output.Imports = CollectExprImports(ctx, opts[0])
	output.Imports = append(output.Imports, valley.Import{
		Path:  "reflect",
		Alias: "reflect",
	})

	output.Code = GenerateStandardConstraint(ctx, predicate, message, details)

	return output, nil
}
