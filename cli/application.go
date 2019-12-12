package cli

import (
	"github.com/seeruk/go-console"
	"github.com/seeruk/valley"
)

// NewApplication returns a new console application. The given constraints are used to produce the
// generate command, allowing binaries with custom validation constraints to be built easily.
func NewApplication(constraints map[string]valley.ConstraintGenerator) *console.Application {
	application := console.NewApplication("valley", "SNAPSHOT")
	application.SetRootCommand(RootCommand(constraints))

	return application
}
