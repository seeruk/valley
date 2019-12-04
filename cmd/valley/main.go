package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/validation"
	"github.com/seeruk/valley/valley"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("valley: a package path, and a config file path must be given")
		os.Exit(1)
	}

	// TODO: This is pretty crap.
	srcPath := os.Args[1]
	configPath := os.Args[2]
	destPath := os.Args[3]

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
	generator := validation.NewGenerator()

	pkg, err := reader.Read(srcPath)
	if err != nil {
		fmt.Printf("valley: failed to read structs in: %q: %v", srcPath, err)
		os.Exit(1)
	}

	buffer, err := generator.Generate(config, pkg)
	if err != nil {
		fmt.Printf("valley: failed to generate code: %v", err)
		os.Exit(1)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("valley: failed to open destination file for writing: %s: %q", destPath, err)
		os.Exit(1)
	}

	// Write output to stdout.
	_, err = io.Copy(destFile, buffer)
	if err != nil {
		fmt.Printf("valley: failed to write generated code: %v", err)
		os.Exit(1)
	}

}
