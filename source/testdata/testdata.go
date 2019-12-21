package testdata

import (
	"github.com/seeruk/valley"

	// Aliased to test that import aliases are also captured.
	c "github.com/seeruk/valley/validation/constraints"
)

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeText string   `json:"some_text"`
	SomeBool bool     `json:"some_bool"`
	SomePtr  *Subject `json:"some_ptr"`
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s Subject) Constraints(t valley.Type) {
	t.Field(s.SomeText).Constraints(c.Required())
	t.Field(s.SomeBool).Constraints(c.Equals(true))
	t.Field(s.SomePtr).Constraints(c.NotNil())
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

// SomeVar is used to ensure generic declarations that aren't types aren't included in the
// information read from Go source files when testing source reading functionality.
var SomeVar = 123

// SomeInterface is used to ensure interfaces aren't included in the information read from Go source
// files when testing source reading functionality.
type SomeInterface interface {
	ThatIsUnused()
}

// SomeFunction is used to ensure function declarations that aren't methods aren't included in the
// information read from Go source files when testing source reading functionality.
func SomeFunction() {
	// No-op.
}
