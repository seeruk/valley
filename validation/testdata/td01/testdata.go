package td01

import (
	"regexp"
	"time"

	"github.com/seeruk/valley"

	// Aliased to test that import aliases are also captured.
	c "github.com/seeruk/valley/validation/constraints"
)

// patternGreeting ...
var patternGreeting = regexp.MustCompile("^Hello")

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeBool  bool              `json:"some_bool"`
	SomeChan  chan struct{}     `json:"some_chan"`
	SomeMap   map[string]int    `json:"some_map"`
	SomePtr   *Subject          `json:"some_ptr"`
	SomeSlice []string          `json:"some_slice"`
	SomeText  string            `json:"some_text"`
	SomeTime  time.Time         `json:"some_time"`
	Secondary *SecondarySubject `json:"secondary"`
}

// Constraints is a valley constraints method used for testing source reading functionality.
func (s Subject) Constraints(t valley.Type) {
	t.Constraints(c.AnyNRequired(3, s.SomeText, s.SomeBool, s.SomePtr, s.SomeSlice, s.SomeMap))
	t.Constraints(c.ExactlyNRequired(2, s.SomeText, s.SomePtr))
	t.Constraints(c.MutuallyExclusive(s.SomeSlice, s.SomeMap))
	t.Constraints(c.MutuallyInclusive(s.SomeText, s.SomePtr))

	t.Field(s.SomeBool).
		Constraints(c.DeepEquals(true), c.Equals(true), c.NotEquals(false))
	t.Field(s.SomeChan).
		Constraints(c.Nil(), c.NotNil())
	t.Field(s.SomeMap).
		Constraints(c.Required(), c.MinLength(1), c.Nil(), c.NotNil()).
		Elements(c.Required(), c.Min(1)).
		Keys(c.Required(), c.MinLength(3))
	t.Field(s.SomePtr).
		Constraints(c.Required(), c.Nil(), c.NotNil())
	t.Field(s.SomeSlice).
		Constraints(c.Required(), c.Length(16), c.MinLength(2), c.MaxLength(128), c.Nil(), c.NotNil()).
		Elements(c.Required(), c.Length(8), c.MinLength(2), c.MaxLength(32))
	t.Field(s.SomeText).
		Constraints(c.Required(), c.Regexp(patternGreeting), c.RegexpString("^Hello")).
		Constraints(c.OneOf("Hello, World!", "Hello, Go!"), c.Predicate(1 == 1, "1 must equal 1"))
	t.Field(s.SomeTime).
		Constraints(c.Required(), c.TimeAfter(time.Now()), c.TimeBefore(time.Now())).
		Constraints(c.TimeStringAfter("2006-01-02"), c.TimeStringBefore("2006-01-02"))

	t.When(s.SomeBool).Field(s.SomePtr).
		Constraints(c.Nil(), c.NotNil())

	t.Field(s.Secondary).Constraints(c.Valid())
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
