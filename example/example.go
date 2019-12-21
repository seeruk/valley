package main

//go:generate valley ./example.go -t json

import (
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/seeruk/valley"
	"github.com/seeruk/valley/validation/constraints"
)

// patternGreeting is a regular expression to test that a string starts with "Hello".
var patternGreeting = regexp.MustCompile("^Hello")

// timeYosemite is a time that represents when Yosemite National Park was founded.
var timeYosemite = time.Date(1890, time.October, 1, 0, 0, 0, 0, time.UTC)

// Example ...
type Example struct {
	Bool      bool                       `json:"bool,omitempty"`
	Chan      <-chan string              `json:"chan" valley:"chan"`
	Text      string                     `json:"text"`
	Texts     []string                   `json:"texts" valley:"texts"`
	TextMap   map[string]string          `json:"text_map"`
	Adults    int                        `json:"adults"`
	Children  int                        `json:"children" valley:"children"`
	Int       int                        `json:"int"`
	Int2      *int                       `json:"int2" valley:"int2"`
	Ints      []int                      `json:"ints"`
	Float     float64                    `json:"float" valley:"float"`
	Time      time.Time                  `json:"time" valley:"time"`
	Times     []time.Time                `json:"times"`
	Nested    *NestedExample             `json:"nested" valley:"nested"`
	Nesteds   []*NestedExample           `json:"nesteds"`
	NestedMap map[NestedExample]struct{} `json:"nested_map" valley:"nested_map"`
}

// Constraints ...
func (e Example) Constraints(t valley.Type) {
	// Constraints on type as a whole.
	t.Constraints(constraints.MutuallyInclusive(e.Text, e.Texts))
	t.Constraints(constraints.MutuallyInclusive(e.Int, e.Int2, e.Ints))
	t.Constraints(constraints.ExactlyNRequired(3, e.Text, e.Int, e.Int2, e.Ints))

	// Field constraints.
	t.Field(e.Bool).
		Constraints(constraints.NotEquals(false), constraints.DeepEquals(true))
	t.Field(e.Chan).
		Constraints(constraints.MaxLength(12))
	t.Field(e.Text).
		Constraints(
			constraints.Required(),
			constraints.Regexp(patternGreeting),
			constraints.MaxLength(12),
			constraints.Length(5),
			constraints.OneOf("Hello, World!", "Hello, SeerUK!", "Hello, GitHub!"),
			constraints.Predicate(
				strings.HasPrefix(e.Text, "custom") && len(e.Text) == 32,
				"value must be a valid custom ID",
			),
		)
	t.Field(e.TextMap).
		Constraints(constraints.Required()).
		Elements(constraints.Required()).
		Keys(constraints.MinLength(10))
	t.Field(e.Int).
		Constraints(constraints.Required())
	t.Field(e.Int2).
		Constraints(constraints.Required(), constraints.NotNil(), constraints.Min(0))
	t.Field(e.Ints).
		Constraints(constraints.Required(), constraints.MaxLength(3)).
		Elements(constraints.Required(), constraints.Min(0))
	t.Field(e.Float).
		Constraints(constraints.Equals(math.Pi))
	t.Field(e.Time).
		Constraints(constraints.TimeBefore(timeYosemite))
	t.Field(e.Times).
		Constraints(constraints.MinLength(1)).
		Elements(constraints.TimeBefore(timeYosemite))
	t.Field(e.Adults).
		Constraints(constraints.Min(1), constraints.Max(9))
	t.Field(e.Children).
		Constraints(
			constraints.Min(0),
			constraints.Equals(e.Adults+2),
			constraints.Max(int(math.Max(float64(8-(e.Adults-1)), 0))),
		)

	// Conditional constraints.
	t.When(len(e.Text) > 32).Field(e.Text).
		Constraints(constraints.Required(), constraints.MinLength(64))

	// Nested constraints to be called.
	t.Field(e.Nested).
		Constraints(constraints.Required(), constraints.Valid())
	t.Field(e.Nesteds).
		Elements(constraints.Valid())
	t.Field(e.NestedMap).
		Keys(constraints.Valid())
}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}

// Constraints ...
func (n NestedExample) Constraints(t valley.Type) {
	t.Field(n.Text).Constraints(constraints.Required())
}
