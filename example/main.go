package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/valley/example/primary"
)

func main() {
	var example primary.Example

	violations := example.Validate()
	spew.Dump(violations)
}
