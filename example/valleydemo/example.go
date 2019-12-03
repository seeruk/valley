//go:generate sh -c "valley ./example.go > example.validate.go"
package valleydemo

// Example ...
type Example struct {
	Text   string        `json:"text"`
	Texts  []string      `json:"texts"`
	Int    int           `json:"int"`
	Ints   []int         `json:"ints"`
	Nested NestedExample `json:"nested"`
}

func (e *Example) Bla() {}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}
