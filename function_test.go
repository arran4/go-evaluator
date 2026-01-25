package evaluator

import (
	"errors"
	"testing"
)

type SumFunc struct{}

func (s SumFunc) Call(args ...interface{}) (interface{}, error) {
	sum := 0.0
	for _, arg := range args {
		n, ok := numeric[float64](arg)
		if !ok {
			return nil, errors.New("invalid argument")
		}
		sum += n
	}
	return sum, nil
}

func TestFunctionExpression(t *testing.T) {
	expr := FunctionExpression{
		Func: SumFunc{},
		Args: []Term{
			Constant{Value: 10},
			Constant{Value: 20},
		},
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 30.0 {
		t.Errorf("expected 30.0, got %v", result)
	}

	// Test nested evaluation
	expr2 := FunctionExpression{
		Func: SumFunc{},
		Args: []Term{
			expr,               // Result of previous sum (30)
			Constant{Value: 5}, // Add 5
		},
	}

	result2, err := expr2.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result2 != 35.0 {
		t.Errorf("expected 35.0, got %v", result2)
	}
}
