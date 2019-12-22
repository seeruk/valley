package source

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/token"
	"io"
	"os/exec"
	"path"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/seeruk/valley"
)

// Read attempts to read a Go file, and based on it's contents return the package name, along
// with an extract of information about the methods and structs in that file.
func Read(fileSet *token.FileSet, file *ast.File, srcPath string) valley.Source {
	var source valley.Source

	modules := readModules()

	for _, imp := range file.Imports {
		impPath := strings.Trim(imp.Path.Value, "\"")

		impName, ok := readModuleName(modules, impPath)
		if !ok {
			// Fallback to "dumb" way of guessing package name based on path.
			impName = path.Base(impPath)
		}

		if imp.Name != nil {
			impName = imp.Name.Name
		}

		source.Imports = append(source.Imports, valley.Import{
			Alias: impName,
			Path:  impPath,
		})
	}

	source.FileName = path.Base(srcPath)
	source.FileSet = fileSet
	source.Package = file.Name.Name
	source.Methods = make(valley.Methods)
	source.Structs = make(valley.Structs)

	if len(file.Decls) > 0 {
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				readFuncDecl(d, &source)
			case *ast.GenDecl:
				readGenDecl(d, &source)
			}
		}
	}

	structNames := make([]string, 0, len(source.Structs))
	for structName := range source.Structs {
		structNames = append(structNames, structName)
	}

	sort.Strings(structNames)

	source.StructNames = structNames

	return source
}

// readModules ...
func readModules() []valley.Module {
	cmd := exec.Command("go", "list", "-m", "-json", "all")

	bs, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: Warn?
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(bs))

	var modules []valley.Module

	for {
		var module valley.Module

		err := decoder.Decode(&module)
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: Warn?
			return nil
		}

		modules = append(modules, module)
	}

	return modules
}

// readModuleName ...
func readModuleName(modules []valley.Module, impPath string) (string, bool) {
	var module valley.Module
	for _, m := range modules {
		if strings.HasPrefix(impPath, m.Path) {
			module = m
		}
	}

	// Didn't find a matching module.
	if module.Path == "" || module.Dir == "" {
		return "", false
	}

	relativePath := strings.TrimPrefix(impPath, module.Path)

	cmd := exec.Command("go", "list", "-json", path.Join(module.Dir, relativePath))

	bs, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: Warn?
		return "", false
	}

	var pkg valley.Package

	err = json.Unmarshal(bs, &pkg)
	if err != nil {
		// TODO: Warn?
		return "", false
	}

	return pkg.Name, true
}

// readFuncDecl reads a Go function declaration and adds contents that are relevant to the given
// valley Source.
func readFuncDecl(d *ast.FuncDecl, source *valley.Source) {
	if d.Recv == nil || len(d.Recv.List) == 0 || len(d.Recv.List[0].Names) == 0 {
		return
	}

	receiver := d.Recv.List[0]
	receiverName := receiver.Names[0].Name
	receiverType := unpackStarExpr(receiver.Type)

	switch t := receiverType.(type) {
	case *ast.Ident:
		source.Methods[t.Name] = append(source.Methods[t.Name], valley.Method{
			Receiver: receiverName,
			Name:     d.Name.Name,
			Params:   d.Type.Params,
			Results:  d.Type.Results,
			Body:     d.Body,
		})
	}
}

// readGenDecl reads a Go generic declaration and adds contents that are relevant to the given
// valley Source.
func readGenDecl(d *ast.GenDecl, source *valley.Source) {
	if d.Tok != token.TYPE {
		return
	}

	for _, spec := range d.Specs {
		// NOTE: Assumed to always succeed because of the token.TYPE check above.
		typeSpec := spec.(*ast.TypeSpec)

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		fields := readStructFields(structType)
		fieldNames := make([]string, 0, len(fields))

		for fieldName := range fields {
			fieldNames = append(fieldNames, fieldName)
		}

		sort.Strings(fieldNames)

		// At this point, we definitely have a struct.
		structName := typeSpec.Name.Name
		source.Structs[structName] = valley.Struct{
			Name:       structName,
			Node:       structType,
			Fields:     fields,
			FieldNames: fieldNames,
		}
	}
}

// readStructFields reads information about the fields on a given struct type, returning them in a
// more easily accessible format, with only the information we need.
func readStructFields(structType *ast.StructType) valley.Fields {
	fields := make(valley.Fields)

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			valleyField := valley.Value{
				Name: name.Name,
				Type: field.Type,
			}

			if field.Tag != nil {
				_, fs := utf8.DecodeRuneInString(field.Tag.Value)
				_, ls := utf8.DecodeLastRuneInString(field.Tag.Value)

				valleyField.Tag = field.Tag.Value[fs : len(field.Tag.Value)-ls]
			}

			fields[name.Name] = valleyField
		}
	}

	return fields
}

// unpackStarExpr ...
// NOTE: This is purposefully _not_ recursive, as it's invalid for a method receiver to be a pointer
// to a pointer to a type, as that would be an unnamed type.
func unpackStarExpr(expr ast.Expr) ast.Expr {
	se, ok := expr.(*ast.StarExpr)
	if !ok {
		return expr
	}

	return se.X
}
