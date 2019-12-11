package main

//go:generate valley ./example.go

import (
	"math"

	"github.com/seeruk/valley/validation/constraints"
	"github.com/seeruk/valley/valley"
)

// Example ...
type Example struct {
	Bool     bool              `json:"bool"`
	Chan     <-chan string     `json:"chan"`
	Text     string            `json:"text"`
	Texts    []string          `json:"texts"`
	TextMap  map[string]string `json:"text_map"`
	Adults   int               `json:"adults"`
	Children int               `json:"children"`
	Int      int               `json:"int"`
	Int2     *int              `json:"int2"`
	Ints     []int             `json:"ints"`
	Nested   *NestedExample    `json:"nested"`
	Nesteds  []*NestedExample  `json:"nesteds"`
}

// Constraints ...
func (e Example) Constraints(t valley.Type) {
	// Constraints on type as a whole.
	t.Constraints(constraints.MutuallyExclusive(e.Text, e.Texts))

	// List of possible constraints to implement:
	// * MutuallyInclusive: If one is set, all of them must be set.
	// * Length: Exactly length of something that can have length calculated on it.
	// * OneOf: Actual value must be equal to one of the given values (maybe tricky?).
	// * AnyNRequired: Similar to MutuallyExclusive, but making at least one of the values be required.
	// * ExactlyNRequired: Similar to MutuallyExclusive, but making exactly one of the values be required.
	// * TimeBefore: Validates that a time is before another.
	// * TimeAfter: Validates that a time is after another.
	// * Pattern: Validates that a string matches the given regular expression.
	//   * Maybe this should add package-local variables for the patterns or something?
	// * Predicate: Custom code... as real code.

	// Example of Predicate constraint.
	//t.Field(e.Text).Constraints(constraints.Predicate(e.Text == "Hello, World!"))

	// Field constraints.
	t.Field(e.Bool).Constraints(constraints.NotEquals(false), constraints.Equals(true))
	t.Field(e.Chan).Constraints(constraints.MaxLength(12))
	t.Field(e.Text).Constraints(constraints.Required())
	t.Field(e.Text).Constraints(constraints.MaxLength(12))
	t.Field(e.TextMap).Constraints(constraints.Required()).
		Elements(constraints.Required())
	t.Field(e.Int).Constraints(constraints.Required())
	t.Field(e.Int2).Constraints(constraints.Required(), constraints.NotNil(), constraints.Min(0))
	t.Field(e.Ints).Constraints(constraints.Required(), constraints.MaxLength(3)).
		Elements(constraints.Required(), constraints.Min(0))
	t.Field(e.Adults).Constraints(constraints.Min(1), constraints.Max(9))
	t.Field(e.Children).Constraints(constraints.Min(0), constraints.Equals(e.Adults+2)).
		Constraints(constraints.Max(int(math.Max(float64(8-(e.Adults-1)), 0))))

	// Nested constraints to be called.
	t.Field(e.Nested).Constraints(constraints.Required(), constraints.Valid())
	t.Field(e.Nesteds).Elements(constraints.Valid())
}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}

// Constraints ...
func (n NestedExample) Constraints(t valley.Type) {
	t.Field(n.Text).Constraints(constraints.Required())
}
