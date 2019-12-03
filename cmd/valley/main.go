package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/seeruk/valley/constraints"

	"github.com/ghodss/yaml"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/valley"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("valley: a config file path must be given")
		os.Exit(1)
	}

	srcPath := os.Args[1]
	configPath := configPathFromSrc(srcPath)

	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("valley: failed to open config file: %q: %v\n", configPath, err)
		os.Exit(1)
	}

	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("valley: failed to read config file: %q: %v\n", configPath, err)
		os.Exit(1)
	}

	var config valley.Config

	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		fmt.Printf("valley: failed to unmarshal config: %q: %v\n", configPath, err)
		os.Exit(1)
	}

	reader := source.NewReader()

	pkgName, structs, err := reader.Read(srcPath)
	if err != nil {
		fmt.Printf("valley: failed to read structs in: %q: %v", srcPath, err)
		os.Exit(1)
	}

	fmt.Printf("package %s\n", pkgName)
	fmt.Println()
	fmt.Println(`import "github.com/seeruk/valley/valley"`)

	for typeName, typ := range config.Types {
		s, ok := structs[typeName]
		if !ok {
			continue
		}

		// TODO: Do this properly... Not sure how.
		firstRune, _ := utf8.DecodeRuneInString(typeName)
		receiver := strings.ToLower(string(firstRune))

		fmt.Println()
		fmt.Printf("func (%s %s) Validate() []valley.ConstraintViolation {\n", receiver, typeName)
		fmt.Println("	var violations []valley.ConstraintViolation")

		for fieldName, fieldConfig := range typ.Fields {
			f, ok := s.Fields[fieldName]
			if !ok {
				// TODO: Error, bad config
				fmt.Printf("valley: field %q does not exist in Go source", fieldName)
				continue
			}

			for _, constraintConfig := range fieldConfig.Constraints {
				constraint, ok := constraints.BuiltIn[constraintConfig.Name]
				if !ok {
					// TODO: Error, bad config.
					fmt.Printf("valley: unknown validation constraint: %q", constraintConfig.Name)
					continue
				}

				value := valley.Value{
					FieldName: fieldName,
					VarName:   fmt.Sprintf("%s.%s", receiver, fieldName),
				}

				code, err := constraint(value, f.Type, constraintConfig.Opts)
				if err != nil {
					// TODO: Error, invalid config.
					fmt.Printf("valley: failed to generate code for %q.%q's %q constraint: %v", typeName, fieldName, constraintConfig.Name, err)
					continue
				}

				fmt.Println(code)
			}
		}

		fmt.Println("	return violations")
		fmt.Printf("}\n")
	}
}

// configPathFromSrc ...
func configPathFromSrc(srcPath string) string {
	// TODO: Is this approach robust enough?
	// TODO: Probably need to be able to specify overrides for this actually. Or can we read the
	// source for a whole package instead? That might work better.
	return strings.TrimSuffix(srcPath, ".go") + ".valley.yml"
}
