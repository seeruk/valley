package validation

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/seeruk/valley"
	"github.com/seeruk/valley/validation/constraints"
)

// Generator is a type used to generate validation code.
type Generator struct {
	constraints   map[string]valley.ConstraintGenerator
	constraintNum int

	cb   *bytes.Buffer
	ipts map[valley.Import]struct{}
	vars map[valley.Variable]struct{}
}

// NewGenerator returns a new Generator instance.
func NewGenerator(constraints map[string]valley.ConstraintGenerator) *Generator {
	return &Generator{
		constraints: constraints,
		cb:          &bytes.Buffer{},
		ipts: map[valley.Import]struct{}{
			{Path: "fmt", Alias: "fmt"}:                         {},
			{Path: "strconv", Alias: "strconv"}:                 {},
			{Path: "github.com/seeruk/valley", Alias: "valley"}: {},
		},
		vars: make(map[valley.Variable]struct{}),
	}
}

// Generate attempts to generate the code (returned as bytes) to validate code in the given package,
// using the given configuration.
func (g *Generator) Generate(config valley.Config, source valley.Source) ([]byte, error) {
	typeNames := make([]string, 0, len(config.Types))
	for typeName := range config.Types {
		typeNames = append(typeNames, typeName)
	}

	// Ensure we generate methods in the same order each time.
	sort.Strings(typeNames)

	for _, typeName := range typeNames {
		err := g.generateType(config, source, typeName)
		if err != nil {
			return nil, err
		}
	}

	buf := &bytes.Buffer{}

	fmt.Fprintln(buf, "// Code generated by valley. DO NOT EDIT.")
	fmt.Fprintf(buf, "package %s\n", source.Package)
	fmt.Fprintln(buf)

	for ipt := range g.ipts {
		if ipt.Alias != "" {
			fmt.Fprintf(buf, "import %s \"%s\"\n", ipt.Alias, ipt.Path)
		} else {
			fmt.Fprintf(buf, "import \"%s\"\n", ipt.Path)
		}
	}

	fmt.Fprintln(buf)
	fmt.Fprintln(buf, "// Reference imports to suppress errors if they aren't otherwise used")
	fmt.Fprintln(buf, "var _ = fmt.Sprintf")
	fmt.Fprintln(buf, "var _ = strconv.Itoa")
	fmt.Fprintln(buf)

	fmt.Fprintln(buf, "// Variables generated by constraints:")
	for v := range g.vars {
		fmt.Fprintf(buf, "var %s = %s\n", v.Name, v.Value)
	}

	fmt.Fprintln(buf)

	_, err := io.Copy(buf, g.cb)
	if err != nil {
		return nil, errors.New("failed to copy code buffer contents to generate buffer")
	}

	return buf.Bytes(), nil
}

// generateType generates the entire Validate method for a particular type (found in the given
// package with the given type name.
func (g *Generator) generateType(config valley.Config, source valley.Source, typeName string) error {
	typ := config.Types[typeName]

	s, ok := source.Structs[typeName]
	if !ok {
		return nil
	}

	// Figure out an "okay" receiver name, based on the first letter of the type.
	firstRune, _ := utf8.DecodeRuneInString(typeName)
	receiver := strings.ToLower(string(firstRune))

	mm, ok := source.Methods[typeName]
	if ok && len(mm) > 0 {
		receiver = mm[0].Receiver
	}

	g.wc("// Validate validates this %s.\n", typeName)
	g.wc("// This method was generated by Valley.\n")
	g.wc("func (%s %s) Validate(path *valley.Path) []valley.ConstraintViolation {\n", receiver, typeName)
	g.wc("	var violations []valley.ConstraintViolation\n")
	g.wc("\n")
	g.wc("	path.Write(\".\")\n\n")

	ctx := valley.Context{
		Source:   source,
		TypeName: typeName,
		Receiver: receiver,
		VarName:  receiver,
	}

	for _, constraint := range typ.Constraints {
		value := valley.Value{
			Name: s.Name,
			Type: s.Node,
		}

		err := g.generateConstraint(ctx, constraint, value)
		if err != nil {
			return err
		}
	}

	fieldNames := make([]string, 0, len(typ.Fields))
	for fieldName := range typ.Fields {
		fieldNames = append(fieldNames, fieldName)
	}

	// Ensure that each field's validation is generated in the same order each time.
	sort.Strings(fieldNames)

	for _, fieldName := range fieldNames {
		fieldConfig := typ.Fields[fieldName]

		f, ok := s.Fields[fieldName]
		if !ok {
			return fmt.Errorf("field %q does not exist in Go source", fieldName)
		}

		ctx.FieldName = fieldName
		ctx.VarName = fmt.Sprintf("%s.%s", receiver, fieldName)
		ctx.Path = fmt.Sprintf("\"%s\"", fieldName)
		ctx.BeforeViolation = fmt.Sprintf("size := path.Write(%s)", ctx.Path)
		ctx.AfterViolation = "path.TruncateRight(size)"

		err := g.generateField(ctx, fieldConfig, f)
		if err != nil {
			return err
		}
	}

	g.wc("	path.TruncateRight(1)\n")
	g.wc("\n")
	g.wc("	return violations\n")
	g.wc("}\n\n")

	return nil
}

// generateField generates all of the code for a specific field.
func (g *Generator) generateField(ctx valley.Context, fieldConfig valley.FieldConfig, value valley.Value) error {
	err := g.generateFieldConstraints(ctx, fieldConfig, value)
	if err != nil {
		return err
	}

	err = g.generateFieldElementsConstraints(ctx, fieldConfig, value)
	if err != nil {
		return err
	}

	g.wc("\n")

	return nil
}

// generateFieldConstraints generates the code for constraints that apply directly to a specific
// field (i.e. not including things like the element constraints).
func (g *Generator) generateFieldConstraints(ctx valley.Context, fieldConfig valley.FieldConfig, value valley.Value) error {
	// Generate the constraint code for the field as a whole.
	for _, constraintConfig := range fieldConfig.Constraints {
		err := g.generateConstraint(ctx, constraintConfig, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateFieldElementsConstraints generate the constraint code for each element of an
// array/map/slice field.
// TODO: Only works with array/slice currently, need more info on Context to handle maps?
func (g *Generator) generateFieldElementsConstraints(ctx valley.Context, fieldConfig valley.FieldConfig, value valley.Value) error {
	if len(fieldConfig.Elements) == 0 {
		return nil
	}

	// TODO: This output might look a bit weird for maps?
	g.wc("	for i, element := range %s {\n", ctx.VarName)

	for _, constraintConfig := range fieldConfig.Elements {
		elementValue := ctx.Clone()
		elementValue.VarName = "element"

		var elementType ast.Expr

		switch t := value.Type.(type) {
		case *ast.ArrayType:
			elementType = t.Elt
			elementValue.Path = fmt.Sprintf("\"%s.[\" + strconv.Itoa(i) + \"]\"", elementValue.FieldName)
		case *ast.MapType:
			// TODO: We don't do anything with keys right now...
			elementType = t.Value
			// TODO: Does this work well for non-string types?
			// NOTE: %%%%v is needed because this is passed to another Sprintf later...
			elementValue.Path = fmt.Sprintf("\"%s.[\" + fmt.Sprintf(\"%%%%v\", i) + \"]\"", elementValue.FieldName)
		default:
			return errors.New("config for elements applied to non-iterable type")
		}

		// Set up the path writing, now we have everything we need.
		elementValue.BeforeViolation = fmt.Sprintf("size := path.Write(%s)", elementValue.Path)

		elementField := valley.Value{
			Name: value.Name,
			Type: elementType,
		}

		err := g.generateConstraint(elementValue, constraintConfig, elementField)
		if err != nil {
			return err
		}
	}

	g.wc("	}\n\n")

	return nil
}

// generateConstraint attempts to generate the code for a particular constraint.
func (g *Generator) generateConstraint(ctx valley.Context, constraintConfig valley.ConstraintConfig, value valley.Value) error {
	constraint, ok := g.constraints[constraintConfig.Name]
	if !ok {
		return fmt.Errorf("unknown validation constraint: %q", constraintConfig.Name)
	}

	g.constraintNum++

	selector := ctx.TypeName
	if ctx.FieldName != "" {
		selector += "." + ctx.FieldName
	}

	pos := ctx.Source.FileSet.Position(constraintConfig.Pos)

	ctx.Constraint = constraintConfig.Name
	ctx.ConstraintNum = g.constraintNum

	output, err := constraint(ctx, value.Type, constraintConfig.Opts)
	switch {
	case errors.Is(err, constraints.ErrTypeWarning):
		// TODO: Need a better way of logging things than this...
		fmt.Printf("valley: warning generating code for %s's %q constraint on line %d, col %d: %v\n",
			selector, constraintConfig.Name, pos.Line, pos.Column, err)
	case err != nil:
		return fmt.Errorf("failed to generate code for %s's %q constraint on line %d, col %d: %v",
			selector, constraintConfig.Name, pos.Line, pos.Column, err)
	}

	for _, ipt := range output.Imports {
		g.ipts[ipt] = struct{}{}
	}

	for _, v := range output.Vars {
		g.vars[v] = struct{}{}
	}

	g.wc(output.Code)
	g.wc("\n")

	return nil
}

// wc writes code to the code buffer.
func (g *Generator) wc(format string, a ...interface{}) {
	fmt.Fprintf(g.cb, format, a...)
}
