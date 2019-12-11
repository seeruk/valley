package main

//go:generate valley ./example.go

import (
	"math"

	"github.com/seeruk/valley/validation/constraints"
	"github.com/seeruk/valley/valley"
)

// Example ...
type Example struct {
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

	// More permissive type-checking in constraints would allow you to use custom types, and Valley
	// wouldn't need to definitely know about them to generate code. Doing this would mean we need
	// some new, more specific constraints:
	// * NotNil: Would specifically work on things like pointers - this one is less of an issue.
	//   * Also, interfaces - which is a more difficult one to read some a source file.
	//   * Also, channels, functions, maps, pointers, slices, and unsafe pointers apparently.
	// * NotEquals: Accepts interface{}, just use whatever is passed, and do a != on the value.
	// * Equals: Accepts interface{}, probably more handy for things like booleans I guess.
	// * Min: Covers numbers really...
	// * MinLength: Covers things that can have a length (including strings).
	// * True: Would cover boolean... but might be better to use Equals(true) / NotEquals(false)
	//
	// These are actually a lot more explicit and obvious compared to Required which is actually a
	// little ambiguous from the outside, and requires more knowledge about the types than we might
	// be able to get from reading the code.

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
	t.Field(e.Children).Constraints(constraints.Min(0)).
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
