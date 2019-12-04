//go:generate sh -c "valley . ./primary.valley.yml > primary_validate.go"
package primary

import "github.com/seeruk/valley/example/secondary"

// Example ...
type Example struct {
	Text      string              `json:"text"`
	Texts     []string            `json:"texts"`
	TextMap   map[string]string   `json:"text_map"`
	Int       int                 `json:"int"`
	Ints      []int               `json:"ints"`
	Nested    NestedExample       `json:"nested"`
	Secondary secondary.Secondary `json:"secondary"`
}

// Bla ...
func (e *Example) Bla() {}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}
