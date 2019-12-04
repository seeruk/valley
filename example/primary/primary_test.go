package primary_test

import (
	"testing"

	"github.com/seeruk/valley/example/primary"
	"github.com/seeruk/valley/valley"
)

func BenchmarkExample_Validate(b *testing.B) {
	var example primary.Example
	var violations []valley.ConstraintViolation

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		violations = example.Validate()
	}

	_ = violations
}
