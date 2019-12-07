//go:generate valley ./example.go
package primary

import (
	"fmt"

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

	fmt.Println("lolwut")

	// Field constraints.
	t.Field(e.Text).Constraints(constraints.Required())
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
