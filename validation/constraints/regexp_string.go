package constraints

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Regexp ...
func RegexpString(regexp string) valley.Constraint {
	return valley.Constraint{}
}

const regexpStringFormat = `
	if %s {
		%s
		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "value must match regular expression",
			Details: map[string]interface{}{
				"regexp": %s.String(),
			},
		})
		%s
	}
`

// regexpStringGenerator ...
func regexpStringGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var predicate string

	if len(opts) != 1 {
		return output, errors.New("expected exactly one option")
	}

	output.Imports = CollectExprImports(ctx, opts[0])
	output.Imports = append(output.Imports, valley.Import{
		Path:  "regexp",
		Alias: "regexp",
	})

	pattern, err := SprintNode(ctx.Source.FileSet, opts[0])
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	patternVarName := GenerateVariableName(ctx)

	// Add the regexp as a variable, that will be compiled when imported - should help performance.
	output.Vars = []valley.Variable{
		{Name: patternVarName, Value: fmt.Sprintf("regexp.MustCompile(%s)", pattern)},
	}

	// Check if the field is a pointer, if so, we'll add a nil check and dereference from there.
	_, isPointer := fieldType.(*ast.StarExpr)

	varName := ctx.VarName
	if isPointer {
		predicate += fmt.Sprintf("%s != nil && ", varName)
		varName = "*" + varName
	}

	predicate += fmt.Sprintf("!%s.MatchString(%s)", patternVarName, varName)

	output.Code = fmt.Sprintf(regexpStringFormat,
		predicate,
		ctx.BeforeViolation,
		patternVarName,
		ctx.AfterViolation,
	)

	return output, regexpStringTypeCheck(fieldType)
}

// regexpStringTypeCheck ...
func regexpStringTypeCheck(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return regexpStringTypeCheck(e.X)
	case *ast.Ident:
		if e.Name == "string" {
			return nil
		}
	}

	return ErrTypeWarning
}
