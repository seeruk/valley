package constraints

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"github.com/seeruk/valley/valley"
)

// BuiltIn is a map of all of the built-in validation constraints provided by Valley. This is
// exposed so that custom code generators can build on the set of built-in rules, and also use the
// logic exposed. It's tricky to otherwise make Valley extensible.
var BuiltIn = map[string]valley.ConstraintGenerator{
	"github.com/seeruk/valley/validation/constraints.Max":               minMax(max),
	"github.com/seeruk/valley/validation/constraints.Min":               minMax(min),
	"github.com/seeruk/valley/validation/constraints.MaxLength":         minMaxLength(maxLength),
	"github.com/seeruk/valley/validation/constraints.MinLength":         minMaxLength(minLength),
	"github.com/seeruk/valley/validation/constraints.MutuallyExclusive": mutuallyExclusive,
	"github.com/seeruk/valley/validation/constraints.Required":          required,
	"github.com/seeruk/valley/validation/constraints.Valid":             valid,
}

// GenerateEmptinessPredicate ...
func GenerateEmptinessPredicate(varName string, fieldType ast.Expr) (string, []valley.Import) {
	switch expr := fieldType.(type) {
	case *ast.StarExpr:
		return fmt.Sprintf("%s == nil", varName), nil
	case *ast.ArrayType, *ast.MapType:
		return fmt.Sprintf("len(%s) == 0", varName), nil
	case *ast.Ident:
		switch expr.Name {
		case "string":
			return fmt.Sprintf("len(%s) == 0", varName), nil
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			return fmt.Sprintf("%s == 0", varName), nil
		}
	}

	// If we can't tell what the type is by reading the source, fall back to reflection in this
	// case. There are more efficient things that a consumer of Valley can do, but this is easy, and
	// also powers MutuallyExclusive, etc.
	return fmt.Sprintf("reflect.ValueOf(%s).IsZero()", varName), []valley.Import{{Path: "reflect"}}
}

// CollectExprImports looks for things that appear to be imports in the given expression, and
// attempts to find matching imports in the given valley.Context, returning any that are found.
func CollectExprImports(ctx valley.Context, expr ast.Expr) []valley.Import {
	var imports []valley.Import

	ast.Inspect(expr, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.SelectorExpr:
			ident, ok := n.X.(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name == ctx.Receiver {
				return true
			}

			for _, imp := range ctx.Source.Imports {
				if imp.Alias == ident.Name {
					imports = append(imports, imp)
				}
			}
		}

		return true
	})

	return imports
}

// SprintNode uses the go/printer package to print an AST node, returning it as a string.
func SprintNode(fileSet *token.FileSet, node ast.Node) (string, error) {
	var buf bytes.Buffer

	err := printer.Fprint(&buf, fileSet, node)
	if err != nil {
		return "", fmt.Errorf("failed to render node: %v", err)
	}

	return buf.String(), nil
}
