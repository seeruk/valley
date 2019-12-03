package main

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/seeruk/valley/constraints"
	"github.com/seeruk/valley/valley"
)

func main() {
	builtIn := constraints.BuiltIn
	_ = builtIn

	field := valley.Value{
		FieldName: "Text",
		VarName:   "r.Text",
	}

	fmt.Println(constraints.Required(field, &ast.Ident{
		Name: "int",
	}, nil))

	wd, _ := os.Getwd()

	configPath := path.Join(wd, "example", "valleydemo", "todo.valley.yml")

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config valley.Config

	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		panic(err)
	}

	spew.Dump(config)

	//os.Open()
	//yaml.Unmarshal()

	//fset := token.NewFileSet()
	//f, err := parser.ParseFile(fset, "src.go", `
	//package example
	//
	//func foo(bar []Bar) {
	//	bar[12].Baz.Qux()
	//}
	//
	//type Example struct {
	//	Text bool
	//	Texts []string
	//	Number int
	//	Numbers []int
	//	TriState *string
	//	TriStates []*string
	//	Object map[string]interface{}
	//	Nested Nested
	//	NestedMaybe *Nested
	//}
	//
	//type Nested struct {
	//	Foo string
	//}
	//	`, 0)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//_ = f
}
