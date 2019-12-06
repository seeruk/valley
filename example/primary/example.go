package primary

// Example ...
type Example struct {
	Text    string            `json:"text"`
	Texts   []string          `json:"texts"`
	TextMap map[string]string `json:"text_map"`
	Int     int               `json:"int"`
	Ints    []int             `json:"ints"`
	Nested  NestedExample     `json:"nested"`
}

// NestedExample ...
type NestedExample struct {
	Text string `json:"text"`
}
