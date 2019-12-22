package config

import (
	"flag"
	"fmt"
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
	tt := []struct {
		name string
		desc string
	}{
		{name: "td01", desc: "should produce valid config from valid source"},
		{name: "td02", desc: "should error if nothing is passed to 'Field'"},
		{name: "td03", desc: "should error if a method is not called on 'Field'"},
		{name: "td04", desc: "should error if the value passed to 'Field' is not a selector"},
		{name: "td05", desc: "should error if the value passed to 'Field' is a selector not on the receiver's type"},
		{name: "td06", desc: "should ignore statements in a constraints method's body that can't be used"},
		{name: "td07", desc: "should error if a value passed to Constraints is not a function call"},
		{name: "td08", desc: "should error if a constraint is used from within the same package"},
		{name: "td09", desc: "should error if a constraint is used from within the same package"},
		{name: "td10", desc: "should error if a selector uses a package that can't be found in the imports list"},
	}

	for _, tc := range tt {
		inFile := fmt.Sprintf("./testdata/%s/testdata.go", tc.name)
		outFile := fmt.Sprintf("./testdata/%s/testdata.txt", tc.name)

		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, inFile, nil, 0)
		require.NoError(t, err)

		src := source.Read(fileSet, file, inFile)

		// Don't output things that will change each run.
		spewer := spew.NewDefaultConfig()
		spewer.DisablePointerAddresses = true
		spewer.SortKeys = true

		config, err := BuildFromSource(src)

		actual := fmt.Sprintf("Description: %s\n\nConfig:\n\n%s\nError:\n\n%s",
			tc.desc,
			spewer.Sdump(config),
			spewer.Sdump(err),
		)

		if *update {
			err := ioutil.WriteFile(outFile, []byte(actual), 0666)
			require.NoError(t, err)
		}

		bs, err := ioutil.ReadFile(outFile)
		require.NoError(t, err)

		assert.Equal(t, string(bs), actual)
	}
}
