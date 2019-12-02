package constraint

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
		// TODO: Calculate Field value properly.
		violations = append(violations, valley.ConstraintViolation{
			Field: path.Render(),
			Message: "a value is required",
		})
	}
`

// Required ...
func Required(value valley.Value, fieldType ast.Expr, _ interface{}) (string, error) {
	var predicate string

	switch expr := fieldType.(type) {
	case *ast.StarExpr:
		predicate = fmt.Sprintf("%s == nil", value.VarName)
	case *ast.ArrayType, *ast.MapType:
		predicate = fmt.Sprintf("%[1]s == nil || len(%[1]s) == 0", value.VarName)
	case *ast.Ident:
		switch expr.Name {
		case "string":
			predicate = fmt.Sprintf("len(%s) == 0", value.VarName)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			predicate = fmt.Sprintf("%s == 0", value.VarName)

		default:
			return "", fmt.Errorf("valley: can't handle %T (%s) in `Required`", fieldType, expr.Name)
		}
	default:
		return "", fmt.Errorf("valley: can't handle %T in `Required`", fieldType)
	}

	return fmt.Sprintf(requiredFormat, predicate, value.FieldName), nil
}
