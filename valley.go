package valley

import (
	"go/ast"
	"go/token"
	"regexp"
	"time"
)

// Built in regular expression patterns.
var (
	PatternUUID = regexp.MustCompile(`^[0-9A-f]{8}-[0-9A-f]{4}-[0-9A-f]{4}-[0-9A-f]{4}-[0-9A-f]{12}$`)
)

// Constraint is used to identify constraints to generate code for in a Go AST.
type Constraint struct{}

// ConstraintGenerator is a function that can generate constraint code.
type ConstraintGenerator func(value Context, fieldType ast.Expr, opts []ast.Expr) (ConstraintGeneratorOutput, error)

// ConstraintGeneratorOutput represents the information needed to write some code segments to a new
// Go file. They can't be written to whilst we're generating code because each constraint could need
// code to be in different parts of the resulting file (e.g. imports).
type ConstraintGeneratorOutput struct {
	Imports []Import
	Vars    []Variable
	Code    string
}

// All possible PathKind values.
const (
	PathKindStruct  PathKind = "struct"
	PathKindField   PathKind = "field"
	PathKindElement PathKind = "element"
	PathKindKey     PathKind = "key"
)

// PathKind enumerates possible path kinds that apply to constraint violations.
type PathKind string

// ConstraintViolation is the result of a validation failure.
type ConstraintViolation struct {
	Path     string                 `json:"path,omitempty"`
	PathKind string                 `json:"path_kind"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// Context is used to inform a ConstraintGenerator about it's environment, mainly to do with which
// part of a type is being validated, and giving important identifiers to ConstraintGenerators.
type Context struct {
	Source    Source
	TypeName  string
	Receiver  string
	FieldName string
	VarName   string
	Path      string
	PathKind  PathKind

	Constraint      string
	ConstraintNum   int
	BeforeViolation string
	AfterViolation  string
}

// Clone returns a clone of this Context by utilising the properties of Go values.
func (c Context) Clone() Context {
	return c
}

// Source represents the information Valley needs about a particular source file.
type Source struct {
	FileName string
	FileSet  *token.FileSet
	Package  string
	Imports  []Import
	Methods  Methods
	Structs  Structs
}

// Import represents information about a Go import that Valley uses to generate code.
type Import struct {
	Path  string
	Alias string
}

// Variable represents information about a Go variable that Valley uses to generate code.
type Variable struct {
	Name  string
	Value string
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

// TimeMustParse ...
func TimeMustParse(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}
