package vali

import "fmt"

// ConstraintViolation ...
type ConstraintViolation struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Constraint ...
type Constraint func(n, t string, opts interface{})

type Field struct {
	Receiver string
	Name     string
}

func (f Field) AsVariable() string {
	return fmt.Sprintf("%s.%s", f.Receiver, f.Name)
}
