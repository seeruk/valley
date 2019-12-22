package validation

import (
	"flag"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/seeruk/valley"
	"github.com/seeruk/valley/config"
	"github.com/seeruk/valley/source"
	"github.com/seeruk/valley/validation/constraints"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden record files?")

func TestNewGenerator(t *testing.T) {
	t.Run("should not return nil", func(t *testing.T) {
		assert.NotNil(t, NewGenerator(constraints.BuiltIn))
	})

	t.Run("should initialise the constraints buffer", func(t *testing.T) {
		assert.NotNil(t, NewGenerator(constraints.BuiltIn).cb)
	})

	t.Run("should create a generator with some imports pre-assigned", func(t *testing.T) {
		generator := NewGenerator(constraints.BuiltIn)

		assert.Len(t, generator.ipts, 3)
		assert.Contains(t, generator.ipts, valley.Import{Path: "fmt", Alias: "fmt"})
		assert.Contains(t, generator.ipts, valley.Import{Path: "strconv", Alias: "strconv"})
		assert.Contains(t, generator.ipts, valley.Import{Path: "github.com/seeruk/valley", Alias: "valley"})
	})

	t.Run("should initialise the vars map", func(t *testing.T) {
		assert.NotNil(t, NewGenerator(constraints.BuiltIn).vars)
	})
}

func TestGenerator_Generate(t *testing.T) {
	tt := []struct {
		name string
		desc string
	}{
		{name: "td01", desc: "should successfully generate code given valid input"},
	}

	for _, tc := range tt {
		inFile := fmt.Sprintf("./testdata/%s/testdata.go", tc.name)
		outFile := fmt.Sprintf("./testdata/%s/testdata.txt", tc.name)

		fileSet := token.NewFileSet()
		file, err := parser.ParseFile(fileSet, inFile, nil, 0)
		require.NoError(t, err)

		src := source.Read(fileSet, file, inFile)
		cfg, err := config.BuildFromSource(src)
		require.NoError(t, err)

		generator := NewGenerator(constraints.BuiltIn)

		bs, err := generator.Generate(cfg, src, "valley")

		// Don't output things that will change each run.
		spewer := spew.NewDefaultConfig()
		spewer.DisablePointerAddresses = true
		spewer.SortKeys = true

		if err == nil {
			bs, err = format.Source(bs)
			require.NoError(t, err)
		}

		actual := fmt.Sprintf("Description: %s\n\nGenerated:\n\n%s\nError:\n\n%s",
			tc.desc,
			string(bs),
			spewer.Sdump(err),
		)

		if *update {
			err := ioutil.WriteFile(outFile, []byte(actual), 0666)
			require.NoError(t, err)
		}

		bs, err = ioutil.ReadFile(outFile)
		require.NoError(t, err)

		assert.Equal(t, string(bs), actual)
	}
}
