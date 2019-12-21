package constraints

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/seeruk/valley"
)

// MutuallyInclusive ...
func MutuallyInclusive(fields ...interface{}) valley.Constraint {
	return valley.Constraint{}
}

// mutuallyInclusiveFormat ...
const mutuallyInclusiveFormat = `
	{
		// MutuallyInclusive uses it's own block to lock down nonEmpty's scope.
		var nonEmpty []string

		%s

		if len(nonEmpty) > 0 && len(nonEmpty) != %d {
			%s
			violations = append(violations, valley.ConstraintViolation{
				Path: path.String(),
				PathKind: %q,
				Message: "fields are mutually inclusive",
				Details: map[string]interface{}{
					"fields": %s,
				},
			})
			%s
		}
	}
`

// mutuallyInclusiveGenerator ...
func mutuallyInclusiveGenerator(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintGeneratorOutput, error) {
	var output valley.ConstraintGeneratorOutput
	var fields []string

	if len(opts) < 2 {
		return output, errors.New("expected at least two options")
	}

	_, ok := fieldType.(*ast.StructType)
	if !ok {
		return output, fmt.Errorf("`MutuallyInclusive` applied to non-struct type")
	}

	for _, opt := range opts {
		pos := ctx.Source.FileSet.Position(opt.Pos())

		selector, ok := opt.(*ast.SelectorExpr)
		if !ok {
			return output, fmt.Errorf("value passed to `MutuallyInclusive` is not a field selector on line %d, col %d", pos.Line, pos.Column)
		}

		selectorOn, ok := selector.X.(*ast.Ident)
		if !ok || selectorOn.Name != ctx.Receiver {
			return output, fmt.Errorf("value passed to `MutuallyInclusive is not a field on receiver type on line %d, col %d", pos.Line, pos.Column)
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

	quotedAliases := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		quotedAliases = append(quotedAliases, fmt.Sprintf("%q", alias))
	}

	fieldDetails := "[]string{"
	fieldDetails += strings.Join(quotedAliases, ", ")
	fieldDetails += "}"

	output.Code = fmt.Sprintf(mutuallyInclusiveFormat,
		strings.Join(predicates, "\n\n"),
		len(opts),
		ctx.BeforeViolation,
		ctx.PathKind,
		fieldDetails,
		ctx.AfterViolation,
	)

	return output, nil
}
