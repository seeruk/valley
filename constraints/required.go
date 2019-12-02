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
func Required(value valley.Value, fieldType ast.Expr, _ interface{}) (string, error) {
	var predicate string

	switch expr := fieldType.(type) {
	case *ast.StarExpr:
		predicate = fmt.Sprintf("%s == nil", value.FullName())
	case *ast.ArrayType, *ast.MapType:
		predicate = fmt.Sprintf("len(%s) == 0", value.FullName())
	case *ast.Ident:
		switch expr.Name {
		case "string":
			predicate = fmt.Sprintf("len(%s) == 0", value.FullName())
		case "int":
			predicate = fmt.Sprintf("%s == 0", value.FullName())
		default:
			return "", fmt.Errorf("valley: can't handle %q in `Required`", fieldType)
		}
	default:
		return "", fmt.Errorf("valley: can't handle %q in `Required`", fieldType)
	}

	return fmt.Sprintf(requiredFormat, predicate, value.Name()), nil
}
