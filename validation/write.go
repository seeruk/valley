package validation

import (
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// FormatAndWrite formats some generated Go source code (the input bs), and writes it to a file.
func FormatAndWrite(bs []byte, destPath string) error {
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to open destination source for writing: %s: %q", destPath, err)
	}

	defer destFile.Close()

	formatted, err := format.Source(bs)
	if err != nil {
		fmt.Println(string(bs))
		return fmt.Errorf("failed to format generated source: %v", err)
	}

	_, err = destFile.Write(formatted)
	if err != nil {
		return fmt.Errorf("failed to write generated source to source: %v", err)
	}

	return nil
}

// FindDestination attempts to find a sensible destination file path. The file path is absolute.
func FindDestination(srcPath string) (string, error) {
	absoluteSrcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return "", err
	}

	directory := filepath.Dir(absoluteSrcPath)
	fileName := filepath.Base(absoluteSrcPath)

	destName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	destName += "_validate.go"

	return path.Join(directory, destName), nil
}
