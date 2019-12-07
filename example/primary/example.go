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
func (e Example) Constraints() {
	// Constraints on type as a whole.
	valley.Constraints(constraints.MutuallyExclusive(e.Text, e.Texts, e.TextMap))

	// Field constraints.
	valley.Field(e.Text).Constraints(constraints.Required())
	valley.Field(e.Int).Constraints(constraints.Required())
	valley.Field(e.Ints).Constraints(constraints.Required()).
		Elements(constraints.Required(), constraints.Min(12))

	valley.Field(e.Nested).Constraints(constraints.Valid())
}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}

// Constraints ...
func (n NestedExample) Constraints() {
	valley.Field(n.Text).Constraints(constraints.Required())
}
