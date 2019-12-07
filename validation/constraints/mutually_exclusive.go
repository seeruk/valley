package constraints

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/seeruk/valley/valley"
)

// mutuallyExclusiveFields ...
type mutuallyExclusiveFields []string

// mutuallyExclusiveFormat ...
const mutuallyExclusiveFormat = `
	{
		// MutuallyExclusive uses it's own block to lock down nonEmpty's scope.  
		var nonEmpty []string

		%s

		if len(nonEmpty) > 1 {
			%s
			violations = append(violations, valley.ConstraintViolation{
				Field:   path.String(),
				Message: "fields are mutually exclusive",
				Details: map[string]interface{}{
					"fields": nonEmpty,
				},
			})
			%s
		}
	}
`

// MutuallyExclusive ...
func MutuallyExclusive(fields ...interface{}) valley.Constraint {
	return valley.Constraint{}
}

// mutuallyExclusive ...
func mutuallyExclusive(ctx valley.Context, fieldType ast.Expr, opts []ast.Expr) (valley.ConstraintOutput, error) {
	var output valley.ConstraintOutput
	var fields mutuallyExclusiveFields

	//err := json.Unmarshal(opts, &fields)
	//if err != nil {
	//	// TODO: Wrap.
	//	return output, err
	//}

	structType, ok := fieldType.(*ast.StructType)
	if !ok {
		return output, fmt.Errorf("`MutuallyExclusive` applied to non-struct type")
	}

	var predicates []string

	// TODO: This is a little gross...
	for _, structField := range structType.Fields.List {
		for _, structFieldName := range structField.Names {
			name := structFieldName.Name
			for _, field := range fields {
				if field == name {
					predicate, err := GenerateEmptinessPredicate(fmt.Sprintf("%s.%s", ctx.VarName, name), structField.Type)
					if err != nil {
						// TODO: Wrap.
						return output, err
					}

					predicates = append(predicates, fmt.Sprintf(`if !(%s) { 
						nonEmpty = append(nonEmpty, "%s") 
					}`, predicate, name))
				}
			}
		}
	}

	output.Code = fmt.Sprintf(mutuallyExclusiveFormat,
		strings.Join(predicates, "\n\n"),
		ctx.BeforeViolation,
		ctx.AfterViolation,
	)

	return output, nil
}
