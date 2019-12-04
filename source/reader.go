package source

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/seeruk/valley/valley"
)

// Reader is a type used to read Go source code, and return information that Valley needs to then
// generate validation code.
type Reader struct {
	// ...
}

// NewReader returns a new Reader instance.
func NewReader() *Reader {
	return &Reader{}
}

// Read attempts to read a Go file, and based on it's contents return the package name, along with
// an extract of information about the structs in that file.
func (r *Reader) Read(fileName string) (valley.Package, error) {
	var pkg valley.Package

	fset := token.NewFileSet()

	// Read the source file with the given name.
	// TODO: We can also read entire directories like this...
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return pkg, err
	}

	if len(f.Decls) > 0 {
		pkg.Structs = make(valley.Structs)

		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
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
				pkg.Structs[structName] = valley.Struct{
					Name:   structName,
					Node:   structType,
					Fields: r.readStructFields(structType),
				}
			}
		}
	}

	return pkg, nil
}

// readStructFields reads information about the fields on a given struct type, returning them in a
// more easily accessible format, with only the information we need.
func (r *Reader) readStructFields(structType *ast.StructType) valley.Fields {
	fields := make(valley.Fields)

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fields[name.Name] = valley.Field{
				Name: name.Name,
				Type: field.Type,
			}
		}
	}

	return fields
}
