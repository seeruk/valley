package valley

// Config ...
type Config struct {
	Types map[string]TypeConfig `json:"types"`
}

// TypeConfig ...
type TypeConfig struct {
	Constraints []interface{}          `json:"constraints"`
	Fields      map[string]FieldConfig `json:"fields"`
}

// FieldConfig ...
type FieldConfig struct {
	Constraints []interface{}  `json:"constraints"`
	Elements    ElementsConfig `json:"elements"`
}

// ElementsConfig ...
type ElementsConfig struct {
	Constraints []interface{} `json:"constraints"`
}
