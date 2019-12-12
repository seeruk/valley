package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Length ...
func Length(length int) valley.Constraint {
	return valley.Constraint{}
}

// MaxLength ...
func MaxLength(max int) valley.Constraint {
	return valley.Constraint{}
}

// MinLength ...
func MinLength(min int) valley.Constraint {
	return valley.Constraint{}
}

// lengthFormat is the format used for rendering a `Min` constraint.
const lengthFormat = `
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

// Possible lengthKind values.
const (
	lengthExact lengthKind = "exactly"
	lengthMax   lengthKind = "maximum"
	lengthMin   lengthKind = "minimum"
)

// lengthKind ...
type lengthKind string

// lengthGenerator ...
func lengthGenerator(kind lengthKind) valley.ConstraintGenerator {
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

		var message, operator string

		switch kind {
		case lengthExact:
			message = "exact length not met"
			operator = "!="
		case lengthMax:
			message = "maximum length exceeded"
			operator = ">"
		case lengthMin:
			message = "minimum length not met"
			operator = "<"
		}

		predicate += fmt.Sprintf("len(%s) %s %s", varName, operator, value)

		output.Code = fmt.Sprintf(lengthFormat,
			predicate,
			ctx.BeforeViolation,
			message,
			kind,
			value,
			ctx.AfterViolation,
		)

		return output, lengthTypeCheck(fieldType)
	}
}

// lengthTypeCheck ...
func lengthTypeCheck(expr ast.Expr) error {
	// This is everything that's supported by reflect.Value.Len() too. The difference here is that
	// this doesn't support custom types that are really any of these allowed types underneath (at
	// least not yet...)
	switch e := expr.(type) {
	case *ast.StarExpr:
		return lengthTypeCheck(e.X)
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
