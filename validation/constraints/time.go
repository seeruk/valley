package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"time"

	"github.com/seeruk/valley"
)

// TimeStringAfter ...
func TimeAfter(after time.Time) valley.Constraint {
	return valley.Constraint{}
}

// TimeStringBefore ...
func TimeBefore(before time.Time) valley.Constraint {
	return valley.Constraint{}
}

// Possible timeStringKind values.
const (
	timeAfter  timeKind = "after"
	timeBefore timeKind = "before"
)

// timeStringKind ...
type timeKind string

// timeStringGenerator ...
func timeGenerator(kind timeKind) valley.ConstraintGenerator {
	return func(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
		var output valley.ConstraintGeneratorOutput
		var predicate, message string

		if len(opts) != 1 {
			return output, errors.New("expected exactly one option")
		}

		timeSelector, err := SprintNode(ctx.Source.FileSet, opts[0])
		if err != nil {
			return output, fmt.Errorf("failed to render expression: %v", err)
		}

		_, isPointer := fieldType.(*ast.StarExpr)

		varName := ctx.VarName
		if isPointer {
			predicate += fmt.Sprintf("%s != nil && ", varName)
			varName = "*" + varName
		}

		// TODO: These messages aren't great - any way to improve them?
		if kind == timeAfter {
			message = "value must be after time"
			predicate += fmt.Sprintf("!%s.After(%s)", varName, timeSelector)
		} else {
			message = "value must be before time"
			predicate += fmt.Sprintf("!%s.Before(%s)", varName, timeSelector)
		}

		details := map[string]interface{}{
			"time": fmt.Sprintf("%s.Format(time.RFC3339)", timeSelector),
		}

		output.Imports = CollectExprImports(ctx, opts[0])
		output.Imports = append(output.Imports, valley.Import{
			Path:  "time",
			Alias: "time",
		})

		output.Code = GenerateStandardConstraint(ctx, predicate, message, details)

		return output, timeTypeCheck(fieldType)
	}
}

// timeStringTypeCheck ...
func timeTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return timeTypeCheck(e)
	case *ast.SelectorExpr:
		return nil
	}

	return ErrTypeWarning
}
