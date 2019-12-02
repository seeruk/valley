package main

import (
	"go/parser"
	"go/token"
)

func main() {
	//field := valley.Field{
	//	Receiver: "r",
	//	Name:     "ProjectID",
	//}
	//
	//fmt.Println(constraints.Required(field, &ast.StarExpr{
	//	X: &ast.Ident{
	//		Name: "int",
	//	},
	//}, nil))

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", `
	package example

	func foo(bar []Bar) {
		bar[12].Baz.Qux()
	}
	
	type Example struct {
		Text string
		Texts []string
		Number int
		Numbers []int
		TriState *string
		TriStates []*string
		Object map[string]interface{}
		Nested Nested
		NestedMaybe *Nested
	}
	
	type Nested struct {
		Foo string
	}
		`, 0)

	if err != nil {
		panic(err)
	}

	_ = f
}
