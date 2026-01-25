package evaluator

import (
	"testing"
)

type benchUser struct {
	Name string
	Age  int
}

func BenchmarkGreaterThanString_ValueString(b *testing.B) {
	u := &benchUser{Name: "charlie"}
	expr := &GreaterThanExpression{Field: "Name", Value: "bob"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr.Evaluate(u)
	}
}

func BenchmarkGreaterThanString_ValueInt(b *testing.B) {
	u := &benchUser{Name: "123"}
	expr := &GreaterThanExpression{Field: "Name", Value: 100}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr.Evaluate(u)
	}
}

func BenchmarkGreaterThanOrEqualString_ValueInt(b *testing.B) {
	u := &benchUser{Name: "123"}
	expr := &GreaterThanOrEqualExpression{Field: "Name", Value: 100}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr.Evaluate(u)
	}
}

func BenchmarkLessThanString_ValueInt(b *testing.B) {
	u := &benchUser{Name: "099"}
	expr := &LessThanExpression{Field: "Name", Value: 100}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr.Evaluate(u)
	}
}
