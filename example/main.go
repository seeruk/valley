package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/valley/example/primary"
)

func main() {
	var example primary.Example

	example.Ints = []int{0, 1, 2, 3, 0, 5}

	violations := example.Validate()
	spew.Dump(violations)
}
