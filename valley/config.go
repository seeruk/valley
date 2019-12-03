package valley

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
	Elements    ElementsConfig     `json:"elements"`
}

// ElementsConfig ...
type ElementsConfig struct {
	Constraints []ConstraintConfig `json:"constraints"`
}

// ConstraintConfig ...
type ConstraintConfig struct {
	Name string      `json:"name"`
	Opts interface{} `json:"opts"`
}
