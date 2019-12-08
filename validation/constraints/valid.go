package constraints

import (
	"bytes"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// Valid ...
func Valid() valley.Constraint {
	return valley.Constraint{}
}

// valid ...
func valid(ctx valley.Context, fieldType ast.Expr, _ []ast.Expr) (valley.ConstraintOutput, error) {
	buf := &bytes.Buffer{}

	_, isPointer := fieldType.(*ast.StarExpr)

	// If we have a pointer to a struct, unpack it and write an if statement.
	if isPointer {
		fmt.Fprintf(buf, "if %s != nil {\n", ctx.VarName)
	}

	fmt.Fprintf(buf, "%s\n", ctx.BeforeViolation)
	fmt.Fprintf(buf, "violations = append(violations, %s.Validate(path)...)\n", ctx.VarName)
	fmt.Fprintf(buf, "%s\n", ctx.AfterViolation)

	// If we have a pointer to a struct, unpack it and write an if statement.
	if isPointer {
		fmt.Fprintf(buf, "}\n")
	}

	return valley.ConstraintOutput{
		Code: buf.String(),
	}, nil
}
