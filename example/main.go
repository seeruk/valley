package main

import (
	"encoding/json"
	"os"

	"github.com/seeruk/valley/valley"
)

func main() {
	var example Example

	//example.Ints = []int{0, 1, 2, 3, 0, 5}

	example.Int = 1
	example.Ints = []int{1}
	example.Text = "text"
	//example.Texts = []string{"text 1", "text 2"}
	//
	example.TextMap = map[string]string{
		"hello": "",
	}

	example.Nested = &NestedExample{}
	example.Nesteds = []*NestedExample{
		{Text: ""},
	}

	violations := example.Validate(valley.NewPath())

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(violations)
}
