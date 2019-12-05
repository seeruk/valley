package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// requiredFormat is the format used for rendering a `Required` constraint.
// TODO: Maybe there could be a before / after value passed in instead?
const requiredFormat = `
	if %s {
		size := path.Write(%s)

		violations = append(violations, valley.ConstraintViolation{
			Field:   path.String(),
			Message: "a value is required",
		})

		path.TruncateRight(size)
	}
`

// Required ...
func Required(value valley.Value, fieldType ast.Expr, _ interface{}) (valley.ConstraintOutput, error) {
	var output valley.ConstraintOutput
	var predicate string

	switch expr := fieldType.(type) {
	case *ast.StarExpr:
		predicate = fmt.Sprintf("%s == nil", value.VarName)
	case *ast.ArrayType, *ast.MapType:
		predicate = fmt.Sprintf("len(%s) == 0", value.VarName)
	case *ast.Ident:
		switch expr.Name {
		case "string":
			predicate = fmt.Sprintf("len(%s) == 0", value.VarName)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			predicate = fmt.Sprintf("%s == 0", value.VarName)
		default:
			return output, fmt.Errorf("valley: can't handle %T (%s) in `Required`", fieldType, expr.Name)
		}
	default:
		return output, fmt.Errorf("valley: can't handle %T in `Required`", fieldType)
	}

	output.Code = fmt.Sprintf(requiredFormat, predicate, value.Path)

	return output, nil
}
