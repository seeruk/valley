package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/valley/source"
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
		fatalf("valley: a package path, and a config file path must be given\n")
	}

	pkg, err := source.NewReader().Read(srcPath)
	if err != nil {
		fatalf("valley: failed to read structs in: %q: %v\n", srcPath, err)
	}

	spew.Dump(pkg)

	//generator := validation.NewGenerator(constraints.BuiltIn)
	//
	//bs, err := generator.Generate(config, pkg)
	//if err != nil {
	//	fatalf("valley: failed to generate code: %v\n", err)
	//}
	//
	//destFile, err := os.Create(destPath)
	//if err != nil {
	//	fatalf("valley: failed to open destination file for writing: %s: %q\n", destPath, err)
	//}
	//
	//defer destFile.Close()
	//
	//formatted, err := format.Source(bs)
	//if err != nil {
	//	fatalf("valley: failed to format generated source: %v\n", err)
	//}
	//
	//_, err = destFile.Write(formatted)
	//if err != nil {
	//	fatalf("valley: failed to write generated source to file: %v\n", err)
	//}
}

// fatalf writes the given formatted message to stdout, then exits the application with an error
// exit code.
func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
