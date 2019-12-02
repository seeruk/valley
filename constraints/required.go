package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley"
)

// Arguments to format are:
//   [1]: Predicate   (e.g. `== ""`)
//   [2]: Field name  (e.g. ProjectID)
const requiredFormat = `
	if %[1]s {
		// TODO: Add prefix to field names.
		violations = append(violations, valley.ConstraintViolation{
			Field: "%[2]s",
			Message: "a value is required",
		})
	}
`

// Required ...
func Required(field valley.Field, fieldType ast.Expr, _ interface{}) (string, error) {
	var predicate string

	switch expr := fieldType.(type) {
	case *ast.StarExpr:
		// Value is a pointer, unpack it:
		predicate = fmt.Sprintf("%s == nil", field.AsVariable())
	case *ast.ArrayType, *ast.MapType:
		predicate = fmt.Sprintf("len(%s) == 0", field.AsVariable())
	case *ast.Ident:
		switch expr.Name {
		case "string":
			predicate = fmt.Sprintf("len(%s) == 0", field.AsVariable())
		case "int":
			predicate = fmt.Sprintf("%s == 0", field.AsVariable())
		default:
			return "", fmt.Errorf("valley: can't handle %q in `Required`", fieldType)
		}
	default:
		return "", fmt.Errorf("valley: can't handle %q in `Required`", fieldType)
	}

	return fmt.Sprintf(requiredFormat, predicate, field.Name), nil
}
