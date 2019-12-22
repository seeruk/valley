package td06

import (
	"fmt"

	"github.com/seeruk/valley"
	"github.com/seeruk/valley/validation/constraints"
)

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeText string `json:"some_text"`
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s Subject) Constraints(t valley.Type) {
	const foo = "bar"

	two := 1 + 1
	two++

	fmt.Println(t, foo, two)

	bar()

	t.Field(s.SomeText).Constraints(constraints.Required())
}

func bar() {}
