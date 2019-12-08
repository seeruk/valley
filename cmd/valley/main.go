package main

import (
	"fmt"
	"go/format"
	"os"

	"github.com/seeruk/valley/config"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/validation"
	"github.com/seeruk/valley/validation/constraints"
)

func main() {
	var srcPath string
	for _, arg := range os.Args[1:] {
		if arg == "--" {
			continue
		}

		srcPath = arg
	}

	if srcPath == "" {
		fatalf("valley: a package path, and a cfg src path must be given\n")
	}

	src, err := source.Read(srcPath)
	if err != nil {
		fatalf("valley: failed to read structs in: %q: %v\n", srcPath, err)
	}

	cfg, err := config.BuildFromSource(src)
	if err != nil {
		fatalf("valley: failed to generate cfg from src: %v\n", err)
	}

	generator := validation.NewGenerator(constraints.BuiltIn)

	bs, err := generator.Generate(cfg, src)
	if err != nil {
		fatalf("valley: failed to generate code: %v\n", err)
	}

	// TODO: Work this out from the input src name.
	destPath := "./example_validate.go"

	destFile, err := os.Create(destPath)
	if err != nil {
		fatalf("valley: failed to open destination src for writing: %s: %q\n", destPath, err)
	}

	defer destFile.Close()

	formatted, err := format.Source(bs)
	if err != nil {
		fatalf("valley: failed to format generated source: %v\n", err)
	}

	_, err = destFile.Write(formatted)
	if err != nil {
		fatalf("valley: failed to write generated source to src: %v\n", err)
	}
}

// fatalf writes the given formatted message to stdout, then exits the application with an error
// exit code.
func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
