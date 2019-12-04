package main

import (
	"fmt"
	"go/format"
	"os"

	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/validation"
	"github.com/seeruk/valley/validation/constraints"
	"github.com/seeruk/valley/valley"
)

func main() {
	if len(os.Args) < 3 {
		fatalf("valley: a package path, and a config file path must be given\n")
	}

	// TODO: This is pretty crap.
	srcPath := os.Args[1]
	configPath := os.Args[2]
	destPath := os.Args[3]

	// TODO: Should the destPath be removed at this point? If it fails to compile, other parts of
	// this process will likely fail. Maybe if a struct field has been changed or something, it
	// might mean that Valley can't run.

	file, err := os.Open(configPath)
	if err != nil {
		fatalf("valley: failed to open config file: %q: %v\n", configPath, err)
	}

	defer file.Close()

	config, err := valley.ReadConfig(file)
	if err != nil {
		fatalf("valley: failed to read config file: %q: %v\n", configPath, err)
	}

	reader := source.NewReader()
	generator := validation.NewGenerator(constraints.BuiltIn)

	pkg, err := reader.Read(srcPath)
	if err != nil {
		fatalf("valley: failed to read structs in: %q: %v\n", srcPath, err)
	}

	bs, err := generator.Generate(config, pkg)
	if err != nil {
		fatalf("valley: failed to generate code: %v\n", err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		fatalf("valley: failed to open destination file for writing: %s: %q\n", destPath, err)
	}

	defer destFile.Close()

	formatted, err := format.Source(bs)
	if err != nil {
		fatalf("valley: failed to format generated source: %v\n", err)
	}

	_, err = destFile.Write(formatted)
	if err != nil {
		fatalf("valley: failed to write generated source to file: %v\n", err)
	}
}

// fatalf writes the given formatted message to stdout, then exits the application with an error
// exit code.
func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
