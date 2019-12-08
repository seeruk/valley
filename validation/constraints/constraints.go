package constraints

import (
	"fmt"
	"go/ast"

	"github.com/seeruk/valley/valley"
)

// BuiltIn is a map of all of the built-in validation constraints provided by Valley. This is
// exposed so that custom code generators can build on the set of built-in rules, and also use the
// logic exposed. It's tricky to otherwise make Valley extensible.
var BuiltIn = map[string]valley.ConstraintGenerator{
	"github.com/seeruk/valley/validation/constraints.Required":          required,
	"github.com/seeruk/valley/validation/constraints.Min":               min,
	"github.com/seeruk/valley/validation/constraints.MutuallyExclusive": mutuallyExclusive,
	"github.com/seeruk/valley/validation/constraints.Valid":             valid,
}

// GenerateEmptinessPredicate ...
func GenerateEmptinessPredicate(varName string, fieldType ast.Expr) (string, error) {
	var predicate string

	switch expr := fieldType.(type) {
	case *ast.StarExpr:
		predicate = fmt.Sprintf("%s == nil", varName)
	case *ast.ArrayType, *ast.MapType:
		predicate = fmt.Sprintf("len(%s) == 0", varName)
	case *ast.Ident:
		switch expr.Name {
		case "string":
			return fmt.Sprintf("len(%s) == 0", varName), nil
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			predicate = fmt.Sprintf("%s == 0", varName)
		default:
			return "", fmt.Errorf("valley: can't generate emptiness predicate for %T (%s)", fieldType, expr.Name)
		}
	default:
		return "", fmt.Errorf("valley: can't generate emptiness predicate for %T", fieldType)
	}

	return predicate, nil
}
