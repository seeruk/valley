package constraints

import "github.com/seeruk/valley/valley"

// BuiltIn is a slice of all of the built-in validation constraints provided by Valley. This is
// exposed so that custom code generators can build on the set of built-in rules, and also use the
// logic exposed. It's tricky to otherwise make Valley extensible.
var BuiltIn = []valley.Constraint{
	Required,
}
