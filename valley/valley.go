package valley

import (
	"go/ast"
)

// ConstraintViolation ...
type ConstraintViolation struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Constraint ...
type Constraint func(value Value, fieldType ast.Expr, opts interface{}) (string, error)

// Value ...
// TODO: Move to a generator package or something? Generator is maybe a poor name, because it makes
// a pretty code type name. Using something like `code` feels a bit crap though too.
type Value struct {
	FieldName string
	VarName   string
}

// Structs is a map from struct name to Struct.
type Structs map[string]Struct

// Struct represents the information we need about a struct in some Go source code.
type Struct struct {
	Name   string
	Node   *ast.StructType
	Fields Fields
}

// Fields is a map from struct field name to Field.
type Fields map[string]Field

// Field represents the information we need about a struct field in some Go source code.
type Field struct {
	Name string
	Type ast.Expr
}
