package valley

import (
	"go/ast"
	"go/token"
)

// Config ...
type Config struct {
	Types map[string]TypeConfig `json:"types"`
}

// TypeConfig ...
type TypeConfig struct {
	Constraints []ConstraintConfig     `json:"constraints"`
	Fields      map[string]FieldConfig `json:"fields"`
}

// FieldConfig ...
type FieldConfig struct {
	Constraints []ConstraintConfig `json:"constraints"`
	Elements    []ConstraintConfig `json:"elements"`
}

// ConstraintConfig ...
type ConstraintConfig struct {
	Name string     `json:"name"`
	Opts []ast.Expr `json:"opts"`
	Pos  token.Pos
}
