package cli

import (
	"errors"
	"fmt"
	"go/format"
	"os"

	"github.com/seeruk/go-console"
	"github.com/seeruk/go-console/parameters"
	"github.com/seeruk/valley/config"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/validation"
	"github.com/seeruk/valley/valley"
)

// RootCommand returns the root console command used when valley is run. This contains the logic to
// orchestrate reading a Go file, building configuration up from that Go file, generating a set of
// validation source code, formatting that source code, and then writing that to a destination file.
func RootCommand(constraints map[string]valley.ConstraintGenerator) *console.Command {
	var srcPath string

	configure := func(def *console.Definition) {
		def.AddArgument(console.ArgumentDefinition{
			Value: parameters.NewStringValue(&srcPath),
			Spec:  "SRC_PATH",
			Desc:  "The path to a file to generate validation code for.",
		})
	}

	execute := func(int *console.Input, output *console.Output) error {
		if srcPath == "" {
			return errors.New("valley: a package path, and a cfg src path must be given")
		}

		src, err := source.Read(srcPath)
		if err != nil {
			return fmt.Errorf("valley: failed to read structs in: %q: %v", srcPath, err)
		}

		cfg, err := config.BuildFromSource(src)
		if err != nil {
			return fmt.Errorf("valley: failed to generate cfg from src: %v", err)
		}

		generator := validation.NewGenerator(constraints)

		bs, err := generator.Generate(cfg, src)
		if err != nil {
			return fmt.Errorf("valley: failed to generate code: %v", err)
		}

		// TODO: Work this out from the input src name.
		destPath := "./example_validate.go"

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("valley: failed to open destination src for writing: %s: %q", destPath, err)
		}

		defer destFile.Close()

		formatted, err := format.Source(bs)
		if err != nil {
			return fmt.Errorf("valley: failed to format generated source: %v", err)
		}

		_, err = destFile.Write(formatted)
		if err != nil {
			return fmt.Errorf("valley: failed to write generated source to src: %v", err)
		}

		return nil
	}

	return &console.Command{
		Description: "Generates validation code by reading a Go file",
		Configure:   configure,
		Execute:     execute,
	}
}
