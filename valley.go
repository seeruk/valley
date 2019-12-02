package valley

import (
	"go/ast"
	"strings"

	"github.com/seeruk/valley/constraint"
)

// DefaultConstraints specifies the list of built-in constraints used by default.
var DefaultConstraints = []Constraint{
	constraint.Required,
}

// ConstraintViolation ...
type ConstraintViolation struct {
	Field   string                 `json:"field"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Constraint ...
type Constraint func(value Value, fieldType ast.Expr, opts interface{}) (string, error)

// Value ...
type Value struct {
	FieldName string
	VarName   string
}

// Path is used to represent the current position in a structure, to output a useful field value to
// identify where a ConstraintViolation occurred.
type Path struct {
	items []string
}

// NewPath returns a new Path instance.
func NewPath() *Path {
	return &Path{}
}

// Push adds an item to the end of the path.
func (r *Path) Push(item string) {
	r.items = append(r.items, item)
}

// Pop removes an item from the end of the path, and returns it.
func (r *Path) Pop() string {
	var p string
	p, r.items = r.items[len(r.items)-1], r.items[:len(r.items)-1]
	return p
}

// Render renders this path as a string, to be sent to the frontend.
func (r *Path) Render() string {
	return strings.Join(r.items, ".")
}
