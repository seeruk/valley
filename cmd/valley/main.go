package main

import (
	"os"

	"github.com/seeruk/valley/cli"
	"github.com/seeruk/valley/validation/constraints"
)

func main() {
	// One-liner main is intended to make it easier to customise Valley, allowing you to implement
	// your own validation constraints. Take a look at some of the built-in constraints. All you
	// need to do is make an application and pass in your set of constraints, like below:
	os.Exit(cli.NewApplication(constraints.BuiltIn).Run(os.Args[1:], os.Environ()))
}
