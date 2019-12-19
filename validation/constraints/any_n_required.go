package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/seeruk/valley"
)

// AnyNRequired ...
func AnyNRequired(n int, fields ...interface{}) valley.Constraint {
	return valley.Constraint{}
}

const anyNRequiredFormat = `
	{
		// AnyNRequired uses it's own block to lock down nonEmpty's scope.
		var nonEmpty []string

		%s

		if len(nonEmpty) < %s {
			%s
			violations = append(violations, valley.ConstraintViolation{
				Path: path.String(),
				PathKind: %q,
				Message: "minimum number of required fields not met",
				Details: map[string]interface{}{
					"num_required": %s,
					"fields": %s,
				},
			})
			%s
		}
	}
`

// anyNRequiredGenerator ...
func anyNRequiredGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var fields []string

	if len(opts) < 3 {
		return output, errors.New("expected at least three options")
	}

	structType, ok := fieldType.(*ast.StructType)
	if !ok {
		return output, fmt.Errorf("`MutuallyExclusive` applied to non-struct type")
	}

	numRequiredExpr := opts[0]
	numRequired, err := SprintNode(ctx.Source.FileSet, numRequiredExpr)
	if err != nil {
		return output, fmt.Errorf("failed to render expression: %v", err)
	}

	for _, opt := range opts[1:] {
		pos := ctx.Source.FileSet.Position(opt.Pos())

		selector, ok := opt.(*ast.SelectorExpr)
		if !ok {
			return output, fmt.Errorf("value passed to `MutuallyExclusive` is not a field selector on line %d, col %d", pos.Line, pos.Column)
		}

		selectorOn, ok := selector.X.(*ast.Ident)
		if !ok || selectorOn.Name != ctx.Receiver {
			return output, fmt.Errorf("value passed to `MutuallyExclusive is not a field on receiver type on line %d, col %d", pos.Line, pos.Column)
		}

		fields = append(fields, selector.Sel.Name)
	}

	var predicates []string

	// TODO: This is a little gross...
	for _, structField := range structType.Fields.List {
		for _, structFieldName := range structField.Names {
			name := structFieldName.Name
			for _, field := range fields {
				if field == name {
					predicate, imports := GenerateEmptinessPredicate(fmt.Sprintf("%s.%s", ctx.VarName, name), structField.Type)
					output.Imports = append(output.Imports, imports...)

					predicates = append(predicates, fmt.Sprintf(`if !(%s) {
						nonEmpty = append(nonEmpty, "%s")
					}`, predicate, name))
				}
			}
		}
	}

	quotedFields := make([]string, 0, len(fields))
	for _, field := range fields {
		quotedFields = append(quotedFields, fmt.Sprintf("%q", field))
	}

	fieldDetails := "[]string{"
	fieldDetails += strings.Join(quotedFields, ", ")
	fieldDetails += "}"

	output.Code = fmt.Sprintf(anyNRequiredFormat,
		strings.Join(predicates, "\n\n"),
		numRequired,
		ctx.BeforeViolation,
		ctx.PathKind,
		numRequired,
		fieldDetails,
		ctx.AfterViolation,
	)

	return output, nil
}
