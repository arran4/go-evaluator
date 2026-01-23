package evaluator

import (
	"testing"
)

func BenchmarkMapAccess(b *testing.B) {
	data := map[string]interface{}{
		"Name": "benchmark",
		"Age":  30,
		"Tags": []string{"a", "b"},
		"Nested": map[string]interface{}{
			"Key": "Value",
		},
	}
	expr := IsExpression{Field: "Name", Value: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !expr.Evaluate(data) {
			b.Fatal("expected true")
		}
	}
}

func BenchmarkMapAccessMiss(b *testing.B) {
	data := map[string]interface{}{
		"Name": "benchmark",
		"Age":  30,
	}
	expr := IsExpression{Field: "Missing", Value: "benchmark"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if expr.Evaluate(data) {
			b.Fatal("expected false")
		}
	}
}
