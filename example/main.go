package main

import (
	"encoding/json"
	"os"

	"github.com/seeruk/valley/example/primary"
)

func main() {
	var example primary.Example

	//example.Ints = []int{0, 1, 2, 3, 0, 5}

	example.Int = 1
	example.Ints = []int{1}
	example.Text = "text"
	//example.Texts = []string{"text 1", "text 2"}
	//
	//example.TextMap = map[string]string{
	//	"Hello": "World!",
	//}

	violations := example.Validate()

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(violations)
}
