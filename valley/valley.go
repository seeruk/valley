package valley

import (
	"go/ast"
	"go/token"
)

// Constraint ...
type Constraint struct{}

// ConstraintGenerator ...
type ConstraintGenerator func(value Context, fieldType ast.Expr, opts []ast.Expr) (ConstraintOutput, error)

// ConstraintOutput ...
type ConstraintOutput struct {
	Imports []Import
	Code    string
}

// ConstraintViolation ...
type ConstraintViolation struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Context ...
// TODO: Move to a generator package or something? Generator is maybe a poor name, because it makes
// a pretty code type name. Using something like `code` feels a bit crap though too.
type Context struct {
	FileSet   *token.FileSet
	TypeName  string
	Receiver  string
	FieldName string
	VarName   string
	Path      string

	BeforeViolation string
	AfterViolation  string
}

// Clone returns a clone of this Context by utilising the properties of Go values.
func (c Context) Clone() Context {
	return c
}

// File ...
type File struct {
	FileSet *token.FileSet
	Package string
	Imports []Import
	Methods Methods
	Structs Structs
}

// Import represents information about a Go import that Valley uses to generate code.
type Import struct {
	Path  string
	Alias string
}

// Methods is a map from struct name to Method.
type Methods map[string][]Method

// Method represents the information we need about a method in some Go source code.
type Method struct {
	Receiver string
	Name     string
	Params   *ast.FieldList
	Results  *ast.FieldList
	Body     *ast.BlockStmt
}

// Structs is a map from struct name to Struct.
type Structs map[string]Struct

// Struct represents the information we need about a struct in some Go source code.
type Struct struct {
	Name   string
	Node   *ast.StructType
	Fields Fields
}

// Fields is a map from struct field name to Value.
type Fields map[string]Value

// Value represents the information we need about a value (e.g. a struct, or a field on a struct) in
// some Go source code.
type Value struct {
	Name string
	Type ast.Expr
}
