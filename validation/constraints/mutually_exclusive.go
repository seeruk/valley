package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/seeruk/valley"
)

// MutuallyExclusive ...
func MutuallyExclusive(fields ...interface{}) valley.Constraint {
	return valley.Constraint{}
}

// mutuallyExclusiveFormat ...
const mutuallyExclusiveFormat = `
	{
		// MutuallyExclusive uses it's own block to lock down nonEmpty's scope.
		var nonEmpty []string

		%s

		if len(nonEmpty) > 1 {
			%s
			violations = append(violations, valley.ConstraintViolation{
				Path: path.String(),
				PathKind: %q,
				Message: "fields are mutually exclusive",
				Details: map[string]interface{}{
					"fields": nonEmpty,
				},
			})
			%s
		}
	}
`

// mutuallyExclusiveGenerator ...
func mutuallyExclusiveGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var fields []string

	if len(opts) < 2 {
		return output, errors.New("expected at least two options")
	}

	_, ok := fieldType.(*ast.StructType)
	if !ok {
		return output, fmt.Errorf("`MutuallyExclusive` applied to non-struct type")
	}

	for _, opt := range opts {
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

	var aliases []string
	var predicates []string

	// TODO: Can this ever fail?
	structType := ctx.Source.Structs[ctx.TypeName]

	// TODO: This is a little gross...
	for _, structFieldName := range structType.FieldNames {
		structField := structType.Fields[structFieldName]

		for _, field := range fields {
			name := structField.Name
			if field != name {
				continue
			}

			alias, err := valley.GetFieldAliasFromTag(name, ctx.TagName, structField.Tag)
			if err != nil {
				return output, fmt.Errorf("failed to generate output field name: %v", err)
			}

			predicate, imports := GenerateEmptinessPredicate(fmt.Sprintf("%s.%s", ctx.VarName, name), structField.Type)
			output.Imports = append(output.Imports, imports...)

			aliases = append(aliases, alias)
			predicates = append(predicates, fmt.Sprintf(`if !(%s) {
				nonEmpty = append(nonEmpty, "%s")
			}`, predicate, alias))
		}
	}

	output.Code = fmt.Sprintf(mutuallyExclusiveFormat,
		strings.Join(predicates, "\n\n"),
		ctx.BeforeViolation,
		ctx.PathKind,
		ctx.AfterViolation,
	)

	return output, nil
}
