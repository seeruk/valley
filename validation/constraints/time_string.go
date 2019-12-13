package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/seeruk/valley"
)

// TimeStringAfter ...
func TimeStringAfter(after string) valley.Constraint {
	return valley.Constraint{}
}

// TimeStringBefore ...
func TimeStringBefore(before string) valley.Constraint {
	return valley.Constraint{}
}

// timeStringFormat is the format used for rendering a `TimeStringAfter` or `TimeStringBefore` constraint.
const timeStringFormat = `
	if %s {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "%s",
			Details: map[string]interface{}{
				"time": %s.Format(time.RFC3339),
			},
		})
		%s
	}
`

// Possible timeStringKind values.
const (
	timeStringAfter  timeStringKind = "after"
	timeStringBefore timeStringKind = "before"
)

// timeStringKind ...
type timeStringKind string

// timeStringGenerator ...
func timeStringGenerator(kind timeStringKind) valley.ConstraintGenerator {
	return func(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
		var output valley.ConstraintGeneratorOutput
		var message string
		var predicate string

		if len(opts) != 1 {
			return output, errors.New("expected exactly one option")
		}

		output.Imports = CollectExprImports(ctx, opts[0])
		output.Imports = append(output.Imports, valley.Import{
			Path:  "time",
			Alias: "time",
		})

		timeString, err := SprintNode(ctx.Source.FileSet, opts[0])
		if err != nil {
			return output, fmt.Errorf("failed to render expression: %v", err)
		}

		timeVarName := GenerateVariableName(ctx)

		output.Vars = []valley.Variable{
			{Name: timeVarName, Value: fmt.Sprintf("valley.TimeMustParse(time.Parse(time.RFC3339, %s))", timeString)},
		}

		_, isPointer := fieldType.(*ast.StarExpr)

		varName := ctx.VarName
		if isPointer {
			predicate += fmt.Sprintf("%s != nil && ", varName)
			varName = "*" + varName
		}

		// TODO: These messages aren't great - any way to improve them?
		if kind == timeStringAfter {
			message = "value must be after time"
			predicate += fmt.Sprintf("!%s.After(%s)", varName, timeVarName)
		} else {
			message = "value must be before time"
			predicate += fmt.Sprintf("!%s.Before(%s)", varName, timeVarName)
		}

		output.Code = fmt.Sprintf(timeStringFormat,
			predicate,
			ctx.BeforeViolation,
			message,
			timeVarName,
			ctx.AfterViolation,
		)

		return output, timeStringTypeCheck(fieldType)
	}
}

// timeStringTypeCheck ...
func timeStringTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return timeStringTypeCheck(e)
	case *ast.SelectorExpr:
		return nil
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			return nil
		}
	}

	return ErrTypeWarning
}
