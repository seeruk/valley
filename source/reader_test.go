package source

import (
	"flag"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden record files?")

func TestRead(t *testing.T) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "./testdata/testdata.go", nil, 0)
	require.NoError(t, err)

	source := Read(fileSet, file, "./testdata/testdata.go")

	t.Run("should set file name on the returned source", func(t *testing.T) {
		assert.Equal(t, "testdata.go", source.FileName)
	})

	t.Run("should set fileset on the returned source", func(t *testing.T) {
		assert.Equal(t, fileSet, source.FileSet)
	})

	t.Run("should set package name on the returned source", func(t *testing.T) {
		assert.Equal(t, "testdata", source.Package)
	})

	t.Run("should set structs on the returned source", func(t *testing.T) {
		require.NotNil(t, source.Structs)
		require.NotNil(t, source.StructNames)

		assert.Len(t, source.Structs, 2)
		assert.Len(t, source.StructNames, 2)
	})

	t.Run("should set methods on the returned source", func(t *testing.T) {
		require.NotNil(t, source.Methods)
		assert.Len(t, source.Methods, 2)
	})

	t.Run("should set imports on the returned source", func(t *testing.T) {
		require.NotNil(t, source.Imports)
		assert.Len(t, source.Imports, 2)
	})

	t.Run("should match the test snapshot", func(t *testing.T) {
		localSource := source
		localSource.FileSet = nil // Already tested.

		// Don't output things that will change each run.
		spewer := spew.NewDefaultConfig()
		spewer.DisablePointerAddresses = true
		spewer.SortKeys = true

		actual := spewer.Sdump(source)
		if *update {
			err := ioutil.WriteFile("./testdata/testdata.txt", []byte(actual), 0666)
			require.NoError(t, err)
		}

		bs, err := ioutil.ReadFile("./testdata/testdata.txt")
		require.NoError(t, err)

		assert.Equal(t, string(bs), actual)
	})
}
