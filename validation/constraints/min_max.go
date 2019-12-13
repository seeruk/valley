package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Max ...
func Max(max int) valley.Constraint {
	return valley.Constraint{}
}

// Min ...
func Min(min int) valley.Constraint {
	return valley.Constraint{}
}

// minMaxFormat is the format used for rendering a `Min` or `Max` constraint.
const minMaxFormat = `
	if %s {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "%s",
			Details: map[string]interface{}{
				"%s": %s,
			},
		})
		%s
	}
`

// Possible minMaxKind values.
const (
	max minMaxKind = "maximum"
	min minMaxKind = "minimum"
)

// minMaxKind ...
type minMaxKind string

// minMaxGenerator ...
func minMaxGenerator(kind minMaxKind) valley.ConstraintGenerator {
	return func(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
		var output valley.ConstraintGeneratorOutput
		var predicate string

		if len(opts) != 1 {
			return output, errors.New("expected exactly one option")
		}

		output.Imports = CollectExprImports(ctx, opts[0])

		// Render the expression passed as an argument to `Min`. We're relying on the fact that the code
		// won't compile if this is configured incorrectly here.
		value, err := SprintNode(ctx.Source.FileSet, opts[0])
		if err != nil {
			return output, fmt.Errorf("failed to render expression: %v", err)
		}

		// Check if the field is a pointer, if so, we'll add a nil check and dereference from there.
		_, isPointer := fieldType.(*ast.StarExpr)

		varName := ctx.VarName
		if isPointer {
			predicate += fmt.Sprintf("%s != nil && ", varName)
			varName = "*" + varName
		}

		message := "maximum value exceeded"
		operator := ">"

		if kind == min {
			message = "minimum value not met"
			operator = "<"
		}

		predicate += fmt.Sprintf("%s %s %s", varName, operator, value)

		output.Code = fmt.Sprintf(minMaxFormat,
			predicate,
			ctx.BeforeViolation,
			message,
			kind,
			value,
			ctx.AfterViolation,
		)

		return output, minMaxTypeCheck(fieldType)
	}
}

// minMaxTypeCheck ...
func minMaxTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return minMaxTypeCheck(e.X)
	case *ast.Ident:
		switch e.Name {
		// TODO: What about... rune, and other built-in types that alias int? For not they'll show
		// a warning I suppose.
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			return nil
		}
	}

	return ErrTypeWarning
}
