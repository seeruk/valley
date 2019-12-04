package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/seeruk/valley/validation"

	"github.com/ghodss/yaml"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/valley"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("valley: a package path, and a config file path must be given")
		os.Exit(1)
	}

	srcPath := os.Args[1]
	configPath := os.Args[2]

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

	buffer := generator.Generate(config, pkg)

	// Write output to stdout.
	_, err = io.Copy(os.Stdout, buffer)
	if err != nil {
		fmt.Printf("valley: failed to write generatede code: %v", err)
		os.Exit(1)
	}
}
