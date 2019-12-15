package constraints

import (
	"go/ast"

	"github.com/seeruk/valley"
)

// Required ...
func Required() valley.Constraint {
	return valley.Constraint{}
}

// requiredGenerator ...
func requiredGenerator(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	predicate, imports := GenerateEmptinessPredicate(ctx.VarName, fieldType)
	return valley.ConstraintGeneratorOutput{
		Imports: imports,
		Code:    GenerateStandardConstraint(ctx, predicate, "a value is required", nil),
	}, nil
}
