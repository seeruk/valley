package valley

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPath(t *testing.T) {
	t.Run("should not return ni", func(t *testing.T) {
		assert.NotNil(t, NewPath())
	})
}

func TestPath_Write(t *testing.T) {
	t.Run("should return the number of bytes written", func(t *testing.T) {
		path := NewPath()
		input := "this is a test"

		expected := len(input)
		actual := path.Write(input)

		assert.Equal(t, expected, actual)
	})

	t.Run("should append the given input to the path", func(t *testing.T) {
		path := NewPath()
		assert.Equal(t, "", path.String())

		path.Write("this is a")
		assert.Equal(t, "this is a", path.String())

		path.Write(" test")
		assert.Equal(t, "this is a test", path.String())
	})
}

func TestPath_TruncateRight(t *testing.T) {
	t.Run("should remove n bytes from the path", func(t *testing.T) {
		path := NewPath()
		input := "this is a test"

		_ = path.Write(input)

		path.TruncateRight(5)
		assert.Equal(t, "this is a", path.String())

		path.TruncateRight(9)
		assert.Equal(t, "", path.String())
	})
}

// NOTE: Path.String() is already well tested enough from the above. Out expectations cover what it
// should be returning.
