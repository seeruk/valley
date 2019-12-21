package valley

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContext_Clone(t *testing.T) {
	t.Run("should return a copy of the given context", func(t *testing.T) {
		context := Context{}

		cloned := context.Clone()
		cloned.TypeName = "test"

		assert.NotEqual(t, context, cloned)
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
