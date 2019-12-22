package td12

import (
	"github.com/seeruk/valley"
	"github.com/seeruk/valley/validation/constraints"
)

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeText string `json:"some_text"`
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s Subject) Constraints(t valley.Type) {
	t.Field(s.SomeText).Constraints(constraints.Required()).
		Constraints("this is invalid")
}
