package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Predicate ...
func Predicate(predicate bool, message string) valley.Constraint {
	return valley.Constraint{}
}

// TODO: How could we get details in here? Maybe even a literal map would be good enough to start?
const predicateFormat = `
	if !(%s) {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: %s,
		})
		%s
	}
`

// predicateGenerator ...
func predicateGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput

	if len(opts) != 2 {
		return output, errors.New("expected exactly two options")
	}

	output.Imports = append(output.Imports, CollectExprImports(ctx, opts[0])...)
	output.Imports = append(output.Imports, CollectExprImports(ctx, opts[1])...)

	predicate, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render predicate expression: %v", err)
	}

	message, err := SprintNode(ctx.Source.FileSet, opts[1])
	if err != nil {
		return output, fmt.Errorf("failed to render message expression: %v", err)
	}

	output.Code = fmt.Sprintf(predicateFormat,
		predicate,
		ctx.BeforeViolation,
		message,
		ctx.AfterViolation,
	)

	return output, nil
}
