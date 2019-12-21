package valley

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_Constraints(t *testing.T) {
	assert.Equal(t, Type{}.Constraints(), Type{})
}

func TestType_Field(t *testing.T) {
	assert.Equal(t, Type{}.Field(nil), Field{})
}

func TestType_When(t *testing.T) {
	assert.Equal(t, Type{}.When(true), Type{})
}

func TestField_Constraints(t *testing.T) {
	assert.Equal(t, Field{}.Constraints(), Field{})
}

func TestField_Elements(t *testing.T) {
	assert.Equal(t, Field{}.Elements(), Field{})
}

func TestField_Keys(t *testing.T) {
	assert.Equal(t, Field{}.Keys(), Field{})
}
