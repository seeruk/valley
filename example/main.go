package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/seeruk/valley"
)

func main() {
	var example Example

	//example.Ints = []int{0, 1, 2, 3, 0, 5}

	example.Adults = 3
	example.Children = 8

	int2 := 12

	example.Bool = true
	example.Int = 5
	example.Int2 = &int2
	example.Ints = []int{1, 2, 3, 4}
	example.Text = "osdfjso i js"
	//example.Texts = []string{"text 1", "text 2"}
	//
	example.TextMap = map[string]string{
		"hello": "",
	}

	example.Nested = &NestedExample{}
	example.Nesteds = []*NestedExample{
		{Text: ""},
	}

	example.Times = []time.Time{
		time.Date(2100, time.October, 01, 0, 0, 0, 0, time.UTC),
	}

	violations := example.Validate(valley.NewPath())

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(violations)
}
