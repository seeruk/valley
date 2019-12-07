//go:generate valley ./example.go
package primary

import (
	"github.com/seeruk/valley/validation/constraints"
	"github.com/seeruk/valley/valley"
)

// Example ...
type Example struct {
	Text    string            `json:"text"`
	Texts   []string          `json:"texts"`
	TextMap map[string]string `json:"text_map"`
	Int     int               `json:"int"`
	Ints    []int             `json:"ints"`
	Nested  NestedExample     `json:"nested"`
}

// Constraints ...
func (e Example) Constraints(t valley.Type) {
	// Constraints on type as a whole.
	t.Constraints(constraints.MutuallyExclusive(e.Text, e.Texts, e.TextMap))

	// List of possible constraints to implement:
	// * MutuallyInclusive: If one is set, all of them must be set.
	// * Min: Min number
	// * Max: Max number
	// * MinLength: Min length of something that can have length calculated on it.
	// * MaxLength: Max length of something that can have length calculated on it.
	// * Length: Exactly length of something that can have length calculated on it.
	// * OneOf: Actual value must be equal to one of the given values (maybe tricky?).
	// * AnyNRequired: Similar to MutuallyExclusive, but making at least one of the values be required.
	// * ExactlyNRequired: Similar to MutuallyExclusive, but making exactly one of the values be required.
	// * TimeBefore: Validates that a time is before another.
	// * TimeAfter: Validates that a time is after another.
	// * Pattern: Validates that a string matches the given regular expression.
	//   * Maybe this should add package-local variables for the patterns or something?
	// * Predicate: Custom code... as real code maybe?

	// Field constraints.
	t.Field(e.Text).Constraints(constraints.Required(), constraints.Predicate(e.Text == "Hello, World!"))
	t.Field(e.Int).Constraints(constraints.Required())
	t.Field(e.Ints).Constraints(constraints.Required()).
		Elements(constraints.Required(), constraints.Min(12))

	// Nested constraints to be called.
	t.Field(e.Nested).Constraints(constraints.Valid())
}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}

// Constraints ...
func (n NestedExample) Constraints(t valley.Type) {
	t.Field(n.Text).Constraints(constraints.Required())
}
