package valley

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Config ...
type Config struct {
	Types map[string]TypeConfig `json:"types"`
}

// ReadConfig attempts to read some Valley configuration from the given Reader.
func ReadConfig(reader io.Reader) (Config, error) {
	var config Config

	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		// TODO: Wrap.
		return config, err
	}

	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		// TODO: Wrap.
		return config, err
	}

	return config, nil
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
	Name string          `json:"name"`
	Opts json.RawMessage `json:"opts"`
}
