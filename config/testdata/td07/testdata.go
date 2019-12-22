package td07

import (
	"github.com/seeruk/valley"
)

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeText string `json:"some_text"`
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s Subject) Constraints(t valley.Type) {
	t.Field(s.SomeText).Constraints("not a function call")
}
