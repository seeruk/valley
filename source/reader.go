package source

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"

	"github.com/seeruk/valley"
)

// Read attempts to read a Go file, and based on it's contents return the package name, along
// with an extract of information about the methods and structs in that file.
func Read(srcPath string) (valley.Source, error) {
	var source valley.Source

	fileSet := token.NewFileSet()

	f, err := parser.ParseFile(fileSet, srcPath, nil, 0)
	if err != nil {
		return source, err
	}

	for _, imp := range f.Imports {
		impPath := strings.Trim(imp.Path.Value, "\"")
		impName := path.Base(impPath)
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
	source.Package = f.Name.Name
	source.Methods = make(valley.Methods)
	source.Structs = make(valley.Structs)

	if len(f.Decls) > 0 {
		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				readFuncDecl(d, &source)
			case *ast.GenDecl:
				readGenDecl(d, &source)
			}
		}
	}

	return source, nil
}

// readFuncDecl reads a Go function declaration and adds contents that are relevant to the given
// valley Source.
func readFuncDecl(d *ast.FuncDecl, source *valley.Source) {
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
		source.Structs[structName] = valley.Struct{
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
