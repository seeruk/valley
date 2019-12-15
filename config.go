package valley

import (
	"go/ast"
	"go/token"
)

// Config represents the configuration for generating an entire set of validation code.
type Config struct {
	Types map[string]TypeConfig `json:"types"`
}

// TypeConfig represents the configuration needed to generate validation code for a specific type.
type TypeConfig struct {
	Constraints []ConstraintConfig     `json:"constraints"`
	Fields      map[string]FieldConfig `json:"fields"`
}

// FieldConfig represents the configuration needed to generate validation code for a specific field
// on a specific type.
type FieldConfig struct {
	Constraints []ConstraintConfig `json:"constraints"`
	Elements    []ConstraintConfig `json:"elements"`
	Keys        []ConstraintConfig `json:"keys"`
}

// ConstraintConfig represents the configuration passed to a ConstraintGenerator to generate some
// code. It's used throughout the configuration structure.
type ConstraintConfig struct {
	Predicate ast.Expr
	Name      string     `json:"name"`
	Opts      []ast.Expr `json:"opts"`
	Pos       token.Pos
}
