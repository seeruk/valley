package td01

import (
	"github.com/seeruk/valley"

	// Aliased to test that import aliases are also captured.
	c "github.com/seeruk/valley/validation/constraints"
)

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeText  string         `json:"some_text"`
	SomeBool  bool           `json:"some_bool"`
	SomePtr   *Subject       `json:"some_ptr"`
	SomeSlice []string       `json:"some_slice"`
	SomeMap   map[string]int `json:"some_map"`
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s Subject) Constraints(t valley.Type) {
	t.Constraints(c.MutuallyExclusive(s.SomeSlice, s.SomeMap))

	t.Field(s.SomeText).
		Constraints(c.Required())
	t.Field(s.SomeBool).
		Constraints(c.Equals(true))
	t.Field(s.SomePtr).
		Constraints(c.NotNil())
	t.Field(s.SomeSlice).
		Constraints(c.MinLength(1), c.MaxLength(128)).
		Elements(c.MinLength(1), c.MaxLength(32))
	t.Field(s.SomeMap).
		Constraints(c.MinLength(1)).
		Elements(c.Min(1)).
		Keys(c.MinLength(3))

	t.When(s.SomeBool).Field(s.SomePtr).
		Constraints(c.NotNil())
}

// SecondarySubject is a type used for testing source reading functionality.
type SecondarySubject struct {
	SomeText string
	SomeBool bool
	SomePtr  *SecondarySubject
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s *SecondarySubject) Constraints(t valley.Type) {
	t.Field(s.SomeText).Constraints(c.Required())
	t.Field(s.SomeBool).Constraints(c.Equals(true))
	t.Field(s.SomePtr).Constraints(c.NotNil())
}
