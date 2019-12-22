package td13

import (
	"go/doc"
)

// Subject is a type used for testing source reading functionality.
type Subject struct {
	SomeText string `json:"some_text"`
}

// NotConstraints1 ...
func (s Subject) NotConstraints1() {}

// NotConstraints2 ...
func (s Subject) NotConstraints2(_ string) {}

// NotConstraints3 ...
func (s Subject) NotConstraints3(_ doc.Type) {}
