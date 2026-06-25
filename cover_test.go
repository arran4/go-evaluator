package evaluator

import (
	"testing"
)

func TestFieldEvaluate(t *testing.T) {
	u := &testUser{Name: "bob", Age: 30}

	f := Field{Name: "Name"}
	val, err := f.Evaluate(u)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if val != "bob" {
		t.Errorf("Expected 'bob', got %v", val)
	}

	f2 := Field{Name: "Missing"}
	_, err = f2.Evaluate(u)
	if err == nil {
		t.Errorf("Expected error for missing field")
	}

	_, err = f.Evaluate(nil)
	if err == nil {
		t.Errorf("Expected error for nil value")
	}
}

func TestSelfEvaluate(t *testing.T) {
	s := Self{}
	val, err := s.Evaluate("test")
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if val != "test" {
		t.Errorf("Expected 'test', got %v", val)
	}
}

func TestBoolTypeEvaluate(t *testing.T) {
	b := BoolType{Term: Constant{Value: true}}
	val, err := b.Evaluate(nil)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if val != true {
		t.Errorf("Expected true, got %v", val)
	}
}

func TestIsTruthy(t *testing.T) {
	truthyCases := []interface{}{
		true,
		"true",
		"1",
		1,
		1.5,
		[]string{"a"},
	}
	for _, c := range truthyCases {
		val, err := IsTruthy(c)
		if err != nil {
			t.Errorf("Expected no error for %v, got %v", c, err)
		}
		if !val {
			t.Errorf("Expected true for %v", c)
		}
	}

	falsyCases := []interface{}{
		false,
		"false",
		"0",
		0,
		0.0,
		nil,
	}
	for _, c := range falsyCases {
		val, err := IsTruthy(c)
		if err != nil {
			t.Errorf("Expected no error for %v, got %v", c, err)
		}
		if val {
			t.Errorf("Expected false for %v", c)
		}
	}
}

func TestIfEvaluate(t *testing.T) {
	i := If{
		Condition: Constant{Value: true},
		Then:      Constant{Value: "yes"},
		Else:      Constant{Value: "no"},
	}
	val, err := i.Evaluate(nil)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if val != "yes" {
		t.Errorf("Expected 'yes', got %v", val)
	}

	i2 := If{
		Condition: Constant{Value: false},
		Then:      Constant{Value: "yes"},
		Else:      Constant{Value: "no"},
	}
	val2, err2 := i2.Evaluate(nil)
	if err2 != nil {
		t.Fatalf("Evaluate error: %v", err2)
	}
	if val2 != "no" {
		t.Errorf("Expected 'no', got %v", val2)
	}
}

func TestComparisonExpressionEvaluate(t *testing.T) {
	cases := []struct {
		op string
		lhs Term
		rhs Term
		expect bool
	}{
		{"eq", Constant{Value: 1}, Constant{Value: 1}, true},
		{"eq", Constant{Value: 1}, Constant{Value: 2}, false},
		{"neq", Constant{Value: 1}, Constant{Value: 2}, true},
		{"neq", Constant{Value: 1}, Constant{Value: 1}, false},
		{"gt", Constant{Value: 2}, Constant{Value: 1}, true},
		{"gt", Constant{Value: 1}, Constant{Value: 2}, false},
		{"gte", Constant{Value: 2}, Constant{Value: 1}, true},
		{"gte", Constant{Value: 2}, Constant{Value: 2}, true},
		{"gte", Constant{Value: 1}, Constant{Value: 2}, false},
		{"lt", Constant{Value: 1}, Constant{Value: 2}, true},
		{"lt", Constant{Value: 2}, Constant{Value: 1}, false},
		{"lte", Constant{Value: 1}, Constant{Value: 2}, true},
		{"lte", Constant{Value: 2}, Constant{Value: 2}, true},
		{"lte", Constant{Value: 2}, Constant{Value: 1}, false},
	}

	for _, c := range cases {
		expr := ComparisonExpression{Operation: c.op, LHS: c.lhs, RHS: c.rhs}
		val, err := expr.Evaluate(nil)
		if err != nil {
			t.Errorf("Error evaluating %v: %v", c, err)
		}
		if val != c.expect {
			t.Errorf("Expected %v for %v, got %v", c.expect, c, val)
		}
	}
}

func TestIContainsExpressionEvaluate(t *testing.T) {
	u := &testUser{Name: "Bob"}

	expr := IContainsExpression{Field: "Name", Value: "ob"}
	val, err := expr.Evaluate(u)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if val != true {
		t.Errorf("Expected true, got %v", val)
	}

	expr2 := IContainsExpression{Field: "Name", Value: "OB"}
	val2, err2 := expr2.Evaluate(u)
	if err2 != nil {
		t.Fatalf("Evaluate error: %v", err2)
	}
	if val2 != true {
		t.Errorf("Expected true, got %v", val2)
	}

	expr3 := IContainsExpression{Field: "Name", Value: "xyz"}
	val3, err3 := expr3.Evaluate(u)
	if err3 != nil {
		t.Fatalf("Evaluate error: %v", err3)
	}
	if val3 != false {
		t.Errorf("Expected false, got %v", val3)
	}

	expr4 := IContainsExpression{Field: "Missing", Value: "a"}
	val4, err4 := expr4.Evaluate(u)
	if err4 != nil {
		t.Fatalf("Evaluate error: %v", err4)
	}
	if val4 != false {
		t.Errorf("Expected false, got %v", val4)
	}

	val5, err5 := expr.Evaluate(nil)
	if err5 != nil {
		t.Fatalf("Evaluate error: %v", err5)
	}
	if val5 != false {
		t.Errorf("Expected false, got %v", val5)
	}
}

type testTypes struct {
	IntVal int
	UintVal uint
	FloatVal float64
	Missing int
}

func TestInequalityExpressionsTypes(t *testing.T) {
	u := &testTypes{IntVal: 10, UintVal: 10, FloatVal: 10.5}

	cases := []struct {
		name   string
		expr   Query
		input  interface{}
		expect bool
	}{
		{"gte int true", Query{Expression: &GreaterThanOrEqualExpression{Field: "IntVal", Value: 10}}, u, true},
		{"gte int false", Query{Expression: &GreaterThanOrEqualExpression{Field: "IntVal", Value: 11}}, u, false},
		{"gte uint true", Query{Expression: &GreaterThanOrEqualExpression{Field: "UintVal", Value: 10}}, u, true},
		{"gte uint false", Query{Expression: &GreaterThanOrEqualExpression{Field: "UintVal", Value: 11}}, u, false},
		{"gte float true", Query{Expression: &GreaterThanOrEqualExpression{Field: "FloatVal", Value: 10.5}}, u, true},
		{"gte float false", Query{Expression: &GreaterThanOrEqualExpression{Field: "FloatVal", Value: 11.5}}, u, false},
		{"lt int true", Query{Expression: &LessThanExpression{Field: "IntVal", Value: 11}}, u, true},
		{"lt int false", Query{Expression: &LessThanExpression{Field: "IntVal", Value: 9}}, u, false},
		{"lt uint true", Query{Expression: &LessThanExpression{Field: "UintVal", Value: 11}}, u, true},
		{"lt uint false", Query{Expression: &LessThanExpression{Field: "UintVal", Value: 9}}, u, false},
		{"lt float true", Query{Expression: &LessThanExpression{Field: "FloatVal", Value: 11.5}}, u, true},
		{"lt float false", Query{Expression: &LessThanExpression{Field: "FloatVal", Value: 9.5}}, u, false},
		{"lte int true", Query{Expression: &LessThanOrEqualExpression{Field: "IntVal", Value: 10}}, u, true},
		{"lte int false", Query{Expression: &LessThanOrEqualExpression{Field: "IntVal", Value: 9}}, u, false},
		{"lte uint true", Query{Expression: &LessThanOrEqualExpression{Field: "UintVal", Value: 10}}, u, true},
		{"lte uint false", Query{Expression: &LessThanOrEqualExpression{Field: "UintVal", Value: 9}}, u, false},
		{"lte float true", Query{Expression: &LessThanOrEqualExpression{Field: "FloatVal", Value: 10.5}}, u, true},
		{"lte float false", Query{Expression: &LessThanOrEqualExpression{Field: "FloatVal", Value: 9.5}}, u, false},
		{"lte float nil false", Query{Expression: &LessThanOrEqualExpression{Field: "FloatVal", Value: 10.5}}, nil, false},
		{"gt int true", Query{Expression: &GreaterThanExpression{Field: "IntVal", Value: 9}}, u, true},
		{"gt uint true", Query{Expression: &GreaterThanExpression{Field: "UintVal", Value: 9}}, u, true},
		{"gt float true", Query{Expression: &GreaterThanExpression{Field: "FloatVal", Value: 9.5}}, u, true},
		{"gt float nil false", Query{Expression: &GreaterThanExpression{Field: "FloatVal", Value: 9.5}}, nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, err := c.expr.Evaluate(c.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if val != c.expect {
				t.Errorf("Expected %v, got %v", c.expect, val)
			}
		})
	}
}

func TestComparisonExpressionEvaluateErrors(t *testing.T) {
	expr := ComparisonExpression{Operation: "eq", LHS: Field{Name: "Missing"}, RHS: Constant{Value: 1}}
	u := &testUser{Name: "Bob"}
	_, err := expr.Evaluate(u)
	if err == nil {
		t.Errorf("Expected error evaluating LHS")
	}

	expr2 := ComparisonExpression{Operation: "eq", LHS: Constant{Value: 1}, RHS: Field{Name: "Missing"}}
	_, err = expr2.Evaluate(u)
	if err == nil {
		t.Errorf("Expected error evaluating RHS")
	}
}

func TestFunctionExpressionEvaluate(t *testing.T) {
	ctx := &Context{
		Functions: map[string]Function{
			"testfunc": dummyFunc{},
		},
	}

	expr := FunctionExpression{Name: "testfunc"}
	val, err := expr.Evaluate(nil, ctx)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if val != "result" {
		t.Errorf("Expected 'result', got %v", val)
	}

	exprMissing := FunctionExpression{Name: "missingfunc"}
	_, err = exprMissing.Evaluate(nil, ctx)
	if err == nil {
		t.Errorf("Expected error for missing function")
	}

	_, err = exprMissing.Evaluate(nil) // no context
	if err == nil {
		t.Errorf("Expected error for missing context")
	}

	exprArgs := FunctionExpression{
		Name: "testfunc",
		Args: []Term{
			Field{Name: "Missing"},
		},
	}
	u := &testUser{Name: "Bob"}
	_, err = exprArgs.Evaluate(u, ctx)
	if err == nil {
		t.Errorf("Expected error evaluating arguments")
	}
}

func TestComparisonExpressionEvaluateBadOp(t *testing.T) {
	expr := ComparisonExpression{Operation: "bad", LHS: Constant{Value: 1}, RHS: Constant{Value: 1}}
	val, err := expr.Evaluate(nil)
	if err != nil {
		t.Errorf("Expected no error for bad op, got %v", err)
	}
	if val != false {
		t.Errorf("Expected false for bad op")
	}
}

type dummyFunc struct{}
func (d dummyFunc) Call(args ...interface{}) (interface{}, error) {
	return "result", nil
}

func TestComparisonExpressionEvaluateContains(t *testing.T) {
	expr := ComparisonExpression{Operation: "contains", LHS: Constant{Value: "hello world"}, RHS: Constant{Value: "world"}}
	val, err := expr.Evaluate(nil)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	if !val {
		t.Errorf("Expected true for contains")
	}

	expr2 := ComparisonExpression{Operation: "contains", LHS: Constant{Value: "hello world"}, RHS: Constant{Value: "xyz"}}
	val2, err2 := expr2.Evaluate(nil)
	if err2 != nil {
		t.Fatalf("Evaluate error: %v", err2)
	}
	if val2 {
		t.Errorf("Expected false for contains")
	}

	expr3 := ComparisonExpression{Operation: "icontains", LHS: Constant{Value: "hello WORLD"}, RHS: Constant{Value: "world"}}
	val3, err3 := expr3.Evaluate(nil)
	if err3 != nil {
		t.Fatalf("Evaluate error: %v", err3)
	}
	if !val3 {
		t.Errorf("Expected true for icontains")
	}

	expr4 := ComparisonExpression{Operation: "icontains", LHS: Constant{Value: "hello WORLD"}, RHS: Constant{Value: "xyz"}}
	val4, err4 := expr4.Evaluate(nil)
	if err4 != nil {
		t.Fatalf("Evaluate error: %v", err4)
	}
	if val4 {
		t.Errorf("Expected false for icontains")
	}
}

func TestComparisonExpressionEvaluateStrings(t *testing.T) {
	u := &testUser{Name: "bob"}

	cases := []struct {
		name   string
		expr   Query
		input  interface{}
		expect bool
	}{
		{"gt true", Query{Expression: &GreaterThanExpression{Field: "Name", Value: "alice"}}, u, true},
		{"gte true", Query{Expression: &GreaterThanOrEqualExpression{Field: "Name", Value: "bob"}}, u, true},
		{"lt true", Query{Expression: &LessThanExpression{Field: "Name", Value: "charlie"}}, u, true},
		{"lte true", Query{Expression: &LessThanOrEqualExpression{Field: "Name", Value: "bob"}}, u, true},
		{"gte false", Query{Expression: &GreaterThanOrEqualExpression{Field: "Name", Value: "charlie"}}, u, false},
		{"lt false", Query{Expression: &LessThanExpression{Field: "Name", Value: "alice"}}, u, false},
		{"lte false", Query{Expression: &LessThanOrEqualExpression{Field: "Name", Value: "alice"}}, u, false},
		{"gte nil false", Query{Expression: &GreaterThanOrEqualExpression{Field: "Name", Value: "bob"}}, nil, false},
		{"lt nil false", Query{Expression: &LessThanExpression{Field: "Name", Value: "charlie"}}, nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, err := c.expr.Evaluate(c.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if val != c.expect {
				t.Errorf("Expected %v, got %v", c.expect, val)
			}
		})
	}
}

func TestComparisonExpressionEvaluateStringsCompare(t *testing.T) {
	u := &testUser{Name: "bob"}

	cases := []struct {
		name   string
		op     string
		field  string
		value  string
		expect bool
	}{
		{"neq true", "neq", "Name", "alice", true},
		{"gt true", "gt", "Name", "alice", true},
		{"gte true", "gte", "Name", "bob", true},
		{"lt true", "lt", "Name", "charlie", true},
		{"lte true", "lte", "Name", "bob", true},
		{"gt false", "gt", "Name", "charlie", false},
		{"gte false", "gte", "Name", "charlie", false},
		{"lt false", "lt", "Name", "alice", false},
		{"lte false", "lte", "Name", "alice", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := ComparisonExpression{Operation: c.op, LHS: Field{Name: c.field}, RHS: Constant{Value: c.value}}
			val, err := expr.Evaluate(u)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if val != c.expect {
				t.Errorf("Expected %v, got %v", c.expect, val)
			}
		})
	}
}
