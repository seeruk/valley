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

const deepEqualsFormat = `
	if !reflect.DeepEqual(%[1]s, %[2]s) {
		%[3]s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "values must be deeply equal",
			Details: map[string]interface{}{
				"deeply_equal_to": %[2]s,
			},
		})
		%[4]s
	}
`

// deepEqualsGenerator ...
func deepEqualsGenerator(ctx valley.Context, _ ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	value, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	output.Imports = CollectExprImports(ctx, opts[0])
	output.Imports = append(output.Imports, valley.Import{
		Path:  "reflect",
		Alias: "reflect",
	})

	output.Code = fmt.Sprintf(deepEqualsFormat,
		ctx.VarName,
		value,
		ctx.BeforeViolation,
		ctx.AfterViolation,
	)

	return output, nil
}
