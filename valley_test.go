package valley

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContext_Clone(t *testing.T) {
	t.Run("should return a copy of the given context", func(t *testing.T) {
		context := Context{}

		cloned := context.Clone()
		cloned.TypeName = "test"

		assert.NotEqual(t, context, cloned)
	})
}

func TestGetFieldAliasFromTag(t *testing.T) {
	t.Run("should return the field name if the tag is empty", func(t *testing.T) {
		alias, err := GetFieldAliasFromTag("testField", "valley", "")
		require.NoError(t, err)
		assert.Equal(t, "testField", alias)
	})

	t.Run("should error if the struct tag is invalid", func(t *testing.T) {
		_, err := GetFieldAliasFromTag("testField", "valley", "this is not a valid tag")
		assert.Error(t, err)
	})

	t.Run("should return the field name if there is no matching tag", func(t *testing.T) {
		alias, err := GetFieldAliasFromTag("testField", "valley", `json:"test"`)
		require.NoError(t, err)
		assert.Equal(t, "testField", alias)
	})

	t.Run("should return the field name if there is an empty matching tag", func(t *testing.T) {
		alias, err := GetFieldAliasFromTag("testField", "valley", `valley:""`)
		require.NoError(t, err)
		assert.Equal(t, "testField", alias)
	})

	t.Run("should return the alias if there is one", func(t *testing.T) {
		alias, err := GetFieldAliasFromTag("testField", "valley", `valley:"test_field"`)
		require.NoError(t, err)
		assert.Equal(t, "test_field", alias)
	})

	t.Run("should return only the section of the tag's value preceding the first comma", func(t *testing.T) {
		alias, err := GetFieldAliasFromTag("testField", "json", `json:"test_field,omitempty"`)
		require.NoError(t, err)
		assert.Equal(t, "test_field", alias)
	})

	t.Run("should remove any excess space from an alias", func(t *testing.T) {
		alias, err := GetFieldAliasFromTag("testField", "valley", `valley:"  test_field  "`)
		require.NoError(t, err)
		assert.Equal(t, "test_field", alias)
	})
}

func TestTimeMustParse(t *testing.T) {
	t.Run("should not panic if the given error is nil", func(t *testing.T) {
		assert.NotPanics(t, func() {
			TimeMustParse(time.Now(), nil)
		})
	})

	t.Run("should panic if the given error is not nil", func(t *testing.T) {
		assert.Panics(t, func() {
			TimeMustParse(time.Parse("2006-01-02", "hello world"))
		})
	})
}
