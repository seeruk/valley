package config

import (
	"flag"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/valley/source"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden record files?")

func TestBuildFromSource(t *testing.T) {
	// TODO: Multiple tests, table driven.
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "./testdata/td01/testdata.go", nil, 0)
	require.NoError(t, err)

	src := source.Read(fileSet, file, "./testdata/td01/testdata.go")

	// Don't output things that will change each run.
	spewer := spew.NewDefaultConfig()
	spewer.DisablePointerAddresses = true
	spewer.SortKeys = true

	config, err := BuildFromSource(src)
	require.NoError(t, err)

	actual := spewer.Sdump(config)
	if *update {
		err := ioutil.WriteFile("./testdata/td01/testdata.txt", []byte(actual), 0666)
		require.NoError(t, err)
	}

	bs, err := ioutil.ReadFile("./testdata/td01/testdata.txt")
	require.NoError(t, err)

	assert.Equal(t, string(bs), actual)
}
