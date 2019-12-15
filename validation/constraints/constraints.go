package constraints

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/seeruk/valley"
)

// BuiltIn is a map of all of the built-in validation constraints provided by Valley. This is
// exposed so that custom code generators can build on the set of built-in rules, and also use the
// logic exposed. It's tricky to otherwise make Valley extensible.
var BuiltIn = map[string]valley.ConstraintGenerator{
	"github.com/seeruk/valley/validation/constraints.DeepEquals":        deepEqualsGenerator,
	"github.com/seeruk/valley/validation/constraints.Equals":            equalsGenerator,
	"github.com/seeruk/valley/validation/constraints.Length":            lengthGenerator(lengthExact),
	"github.com/seeruk/valley/validation/constraints.Max":               minMaxGenerator(max),
	"github.com/seeruk/valley/validation/constraints.MaxLength":         lengthGenerator(lengthMax),
	"github.com/seeruk/valley/validation/constraints.Min":               minMaxGenerator(min),
	"github.com/seeruk/valley/validation/constraints.MinLength":         lengthGenerator(lengthMin),
	"github.com/seeruk/valley/validation/constraints.MutuallyExclusive": mutuallyExclusiveGenerator,
	"github.com/seeruk/valley/validation/constraints.MutuallyInclusive": mutuallyInclusiveGenerator,
	"github.com/seeruk/valley/validation/constraints.NotEquals":         notEqualsGenerator,
	"github.com/seeruk/valley/validation/constraints.NotNil":            notNilGenerator,
	"github.com/seeruk/valley/validation/constraints.Predicate":         predicateGenerator,
	"github.com/seeruk/valley/validation/constraints.Regexp":            regexpGenerator,
	"github.com/seeruk/valley/validation/constraints.RegexpString":      regexpStringGenerator,
	"github.com/seeruk/valley/validation/constraints.Required":          requiredGenerator,
	"github.com/seeruk/valley/validation/constraints.TimeAfter":         timeGenerator(timeAfter),
	"github.com/seeruk/valley/validation/constraints.TimeBefore":        timeGenerator(timeBefore),
	"github.com/seeruk/valley/validation/constraints.TimeStringAfter":   timeStringGenerator(timeStringAfter),
	"github.com/seeruk/valley/validation/constraints.TimeStringBefore":  timeStringGenerator(timeStringBefore),
	"github.com/seeruk/valley/validation/constraints.Valid":             validGenerator,
}

var (
	// ErrTypeWarning is an error returned when a constraint might have been used on an unsupported
	// type. Some constraints may choose to be permissive and continue anyway. This error will only
	// result in a warning being printed. It's up to the constraint if it halts constraint
	// generation (i.e. if that constraint will produce no code, but execution will continue...)
	ErrTypeWarning = errors.New("type used may not produce valid code (is it a custom type?)")
)

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

// GenerateStandardConstraint ...
func GenerateStandardConstraint(ctx valley.Context, predicate, message string, details map[string]interface{}) string {
	constraintFormat := `
		if %s {
			%s
			violations = append(violations, valley.ConstraintViolation{
				Path: path.String(),
				PathKind: %q,
				Message: %q,
				%s
			})
			%s
		}
	`

	var detailsCode string
	if len(details) > 0 {
		detailsCode += "Details: map[string]interface{}{\n"
		for k, v := range details {
			detailsCode += fmt.Sprintf("%q: %v,\n", k, v)
		}
		detailsCode += "},\n"
	}

	return fmt.Sprintf(constraintFormat,
		predicate,
		ctx.BeforeViolation,
		ctx.PathKind,
		message,
		detailsCode,
		ctx.AfterViolation,
	)
}

// GenerateVariableName ...
func GenerateVariableName(ctx valley.Context) string {
	re := regexp.MustCompile(`([^A-z0-9])`)

	return fmt.Sprintf("%s_%s_%d",
		lcfirst(re.ReplaceAllString(ctx.Constraint, "_")),
		ucfirst(strings.TrimSuffix(ctx.Source.FileName, filepath.Ext(ctx.Source.FileName))),
		ctx.ConstraintNum,
	)
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

func lcfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		return u + str[len(u):]
	}
	return ""
}

func ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}
