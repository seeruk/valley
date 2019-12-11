package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// MaxLength ...
func MaxLength(max int) valley.Constraint {
	return valley.Constraint{}
}

// MinLength ...
func MinLength(min int) valley.Constraint {
	return valley.Constraint{}
}

// minMaxFormat is the format used for rendering a `Min` constraint.
const minFormat = `
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
	maxLength minMaxLengthKind = "maximum"
	minLength minMaxLengthKind = "minimum"
)

// minMaxKind ...
type minMaxLengthKind string

// minMax ...
func minMaxLength(kind minMaxLengthKind) valley.ConstraintGenerator {
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

		// TODO: Check the type of the field underneath, it needs to be a string/slice/array/map/channel(?) still.

		varName := ctx.VarName
		if isPointer {
			predicate += fmt.Sprintf("%s != nil && ", varName)
			// Herein we'll be using the de-referenced value.
			varName = "*" + varName
		}

		message := "maximum length exceeded"
		operator := ">"

		if kind == minLength {
			message = "minimum length not met"
			operator = "<"
		}

		predicate += fmt.Sprintf("len(%s) %s %s", varName, operator, value)

		output.Code = fmt.Sprintf(minFormat,
			predicate,
			ctx.BeforeViolation,
			message,
			kind,
			value,
			ctx.AfterViolation,
		)

		return output, minMaxLengthTypeCheck(fieldType)
	}
}

// minMaxLengthTypeCheck ...
func minMaxLengthTypeCheck(expr ast.Expr) error {
	// This is everything that's supported by reflect.Value.Len() too. The difference here is that
	// this doesn't support custom types that are really any of these allowed types underneath (at
	// least not yet...)
	switch e := expr.(type) {
	case *ast.StarExpr:
		return minMaxLengthTypeCheck(e.X)
	case *ast.ArrayType:
		return nil
	case *ast.ChanType:
		return nil
	case *ast.MapType:
		return nil
	case *ast.Ident:
		if e.Name == "string" {
			return nil
		}
	}

	return ErrTypeWarning
}