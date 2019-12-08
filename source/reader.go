package source

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/seeruk/valley/valley"
)

// Read attempts to read a Go file, and based on it's contents return the package name, along
// with an extract of information about the methods and structs in that file.
func Read(srcPath string) (valley.File, error) {
	var file valley.File

	fileSet := token.NewFileSet()

	f, err := parser.ParseFile(fileSet, srcPath, nil, 0)
	if err != nil {
		// TODO: Wrap.
		return file, err
	}

	for _, imp := range f.Imports {
		impPath := strings.Trim(imp.Path.Value, "\"")
		impName := path.Base(impPath)
		if imp.Name != nil {
			impName = imp.Name.Name
		}

		file.Imports = append(file.Imports, valley.Import{
			Alias: impName,
			Path:  impPath,
		})
	}

	file.Name = path.Base(srcPath)
	file.FileSet = fileSet
	file.Package = f.Name.Name
	file.Methods = make(valley.Methods)
	file.Structs = make(valley.Structs)

	if len(f.Decls) > 0 {
		for _, decl := range f.Decls {
			// TODO: Split into multiple methods.
			switch d := decl.(type) {
			case *ast.FuncDecl:
				readFuncDecl(d, &file)
			case *ast.GenDecl:
				readGenDecl(d, &file)
			}
		}
	}

	return file, nil
}

// readFuncDecl reads a Go function declaration and adds contents that are relevant to the given
// valley File.
func readFuncDecl(d *ast.FuncDecl, file *valley.File) {
	if d.Recv == nil {
		return
	}

	// This should probably never happen.
	if len(d.Recv.List) != 1 {
		return
	}

	receiver := d.Recv.List[0]
	receiverName := receiver.Names[0].Name
	receiverType := unpackStarExpr(receiver.Type)

	switch t := receiverType.(type) {
	case *ast.Ident:
		file.Methods[t.Name] = append(file.Methods[t.Name], valley.Method{
			Receiver: receiverName,
			Name:     d.Name.Name,
			Params:   d.Type.Params,
			Results:  d.Type.Results,
			Body:     d.Body,
		})
	}
}

// readGenDecl reads a Go generic declaration and adds contents that are relevant to the given
// valley File.
func readGenDecl(d *ast.GenDecl, file *valley.File) {
	if d.Tok != token.TYPE {
		return
	}

	for _, spec := range d.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// At this point, we definitely have a struct.
		structName := typeSpec.Name.Name
		file.Structs[structName] = valley.Struct{
			Name:   structName,
			Node:   structType,
			Fields: readStructFields(structType),
		}
	}
}

// readStructFields reads information about the fields on a given struct type, returning them in a
// more easily accessible format, with only the information we need.
func readStructFields(structType *ast.StructType) valley.Fields {
	fields := make(valley.Fields)

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fields[name.Name] = valley.Value{
				Name: name.Name,
				Type: field.Type,
			}
		}
	}

	return fields
}

// unpackStarExpr ...
func unpackStarExpr(expr ast.Expr) ast.Expr {
	se, ok := expr.(*ast.StarExpr)
	if !ok {
		return expr
	}

	result := se.X
	if se, ok = result.(*ast.StarExpr); ok {
		return unpackStarExpr(se)
	}

	return result
}
