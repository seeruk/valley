package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/seeruk/valley"
)

// OneOf ...
func OneOf(values ...interface{}) valley.Constraint {
	return valley.Constraint{}
}

// oneOfGenerator ...
func oneOfGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	if len(opts) < 2 {
		return output, errors.New("expected at least two options")
	}

	var predicates []string

	for _, opt := range opts {
		output.Imports = append(output.Imports, CollectExprImports(ctx, opt)...)

		value, err := SprintNode(ctx.Source.FileSet, opt)
		if err != nil {
			return output, fmt.Errorf("failed to render expression: %v", err)
		}

		predicates = append(predicates, fmt.Sprintf("%s != %s", ctx.VarName, value))
	}

	output.Code = GenerateStandardConstraint(ctx,
		strings.Join(predicates, " && "),
		"value must be one of the allowed values",
		nil,
	)

	return output, nil
}
