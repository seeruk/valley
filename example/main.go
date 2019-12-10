package main

import (
	"encoding/json"
	"os"

	"github.com/seeruk/valley/valley"
)

func main() {
	var example Example

	//example.Ints = []int{0, 1, 2, 3, 0, 5}

	example.Adults = 3
	example.Children = 8

	example.Int = 5
	example.Int2 = 12
	example.Ints = []int{1, 2, 3, 4}
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
