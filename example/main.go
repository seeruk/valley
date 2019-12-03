package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/valley/example/valleydemo"
)

func main() {
	var example valleydemo.Example

	violations := example.Validate()
	spew.Dump(violations)
}
