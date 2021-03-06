package cli

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/seeruk/go-console"
	"github.com/seeruk/go-console/parameters"
	"github.com/seeruk/valley"
	"github.com/seeruk/valley/config"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/validation"
)

// RootCommand returns the root console command used when valley is run. This contains the logic to
// orchestrate reading a Go file, building configuration up from that Go file, generating a set of
// validation source code, formatting that source code, and then writing that to a destination file.
func RootCommand(constraints map[string]valley.ConstraintGenerator) *console.Command {
	var srcPath string
	var destPath string

	tagName := "valley"

	configure := func(def *console.Definition) {
		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&destPath),
			Spec:  "-o,--output=DEST",
			Desc:  "Write output to DEST instead of the default '_validate.go'",
		})

		def.AddOption(console.OptionDefinition{
			Value: parameters.NewStringValue(&tagName),
			Spec:  "-t,--tag=NAME",
			Desc:  "Use the given tag name to override field names in generated output (Default: 'valley')",
		})

		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&srcPath),
			Spec:  "SOURCE",
			Desc:  "The path to a file to generate validation code for.",
		})
	}

	execute := func(int *console.Input, output *console.Output) error {
		fileSet := token.NewFileSet()

		file, err := parser.ParseFile(fileSet, srcPath, nil, 0)
		if err != nil {
			return fmt.Errorf("failed to parse source: %v", err)
		}

		src := source.Read(fileSet, file, srcPath)

		cfg, err := config.BuildFromSource(src)
		if err != nil {
			return fmt.Errorf("failed to generate config from source: %v", err)
		}

		generator := validation.NewGenerator(constraints)

		bs, err := generator.Generate(cfg, src, tagName)
		if err != nil {
			return fmt.Errorf("failed to generate validation code: %v", err)
		}

		if destPath == "" {
			destPath = validation.FindDestination(srcPath)
		}

		err = validation.FormatAndWrite(bs, destPath)
		if err != nil {
			return fmt.Errorf("failed to write generated code to destination file: %v", err)
		}

		return nil
	}

	return &console.Command{
		Description: "Generates validation code by reading a Go file",
		Configure:   configure,
		Execute:     execute,
	}
}
