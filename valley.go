package valley

import "go/ast"

// ConstraintViolation ...
type ConstraintViolation struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Constraint ...
type Constraint func(value Value, fieldType ast.Expr, opts interface{}) (string, error)

type Value interface {
	isValue()
}

type Identifier struct {
	Name string
}

type Selector struct {
	Name string
	On   Value
}
