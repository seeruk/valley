package validation

import (
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FormatAndWrite formats some generated Go source code (the input bs), and writes it to a file.
func FormatAndWrite(bs []byte, destPath string) error {
	destFile, err := createFile(destPath)
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
func FindDestination(srcPath string) string {
	destName := strings.TrimSuffix(srcPath, filepath.Ext(srcPath))
	destName += "_validate.go"

	return destName
}

// Testing helpers.

// create is used to provide an easier interface to use in tests to avoid creating real files.
var createFile = func(name string) (io.WriteCloser, error) {
	return os.Create(name)
}
