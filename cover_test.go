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

	exprGteInt := GreaterThanOrEqualExpression{Field: "IntVal", Value: 10}
	val, _ := exprGteInt.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprGteInt = GreaterThanOrEqualExpression{Field: "IntVal", Value: 11}
	val, _ = exprGteInt.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprGteUint := GreaterThanOrEqualExpression{Field: "UintVal", Value: 10}
	val, _ = exprGteUint.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprGteUint = GreaterThanOrEqualExpression{Field: "UintVal", Value: 11}
	val, _ = exprGteUint.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprGteFloat := GreaterThanOrEqualExpression{Field: "FloatVal", Value: 10.5}
	val, _ = exprGteFloat.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprGteFloat = GreaterThanOrEqualExpression{Field: "FloatVal", Value: 11.5}
	val, _ = exprGteFloat.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprLtInt := LessThanExpression{Field: "IntVal", Value: 11}
	val, _ = exprLtInt.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprLtInt = LessThanExpression{Field: "IntVal", Value: 9}
	val, _ = exprLtInt.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprLtUint := LessThanExpression{Field: "UintVal", Value: 11}
	val, _ = exprLtUint.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprLtUint = LessThanExpression{Field: "UintVal", Value: 9}
	val, _ = exprLtUint.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprLtFloat := LessThanExpression{Field: "FloatVal", Value: 11.5}
	val, _ = exprLtFloat.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprLtFloat = LessThanExpression{Field: "FloatVal", Value: 9.5}
	val, _ = exprLtFloat.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprLteInt := LessThanOrEqualExpression{Field: "IntVal", Value: 10}
	val, _ = exprLteInt.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprLteInt = LessThanOrEqualExpression{Field: "IntVal", Value: 9}
	val, _ = exprLteInt.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprLteUint := LessThanOrEqualExpression{Field: "UintVal", Value: 10}
	val, _ = exprLteUint.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprLteUint = LessThanOrEqualExpression{Field: "UintVal", Value: 9}
	val, _ = exprLteUint.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	exprLteFloat := LessThanOrEqualExpression{Field: "FloatVal", Value: 10.5}
	val, _ = exprLteFloat.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprLteFloat = LessThanOrEqualExpression{Field: "FloatVal", Value: 9.5}
	val, _ = exprLteFloat.Evaluate(u)
	if val {
		t.Errorf("Expected false")
	}

	val, _ = exprLteFloat.Evaluate(nil)
	if val {
		t.Errorf("Expected false")
	}

	exprGtInt := GreaterThanExpression{Field: "IntVal", Value: 9}
	val, _ = exprGtInt.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprGtUint := GreaterThanExpression{Field: "UintVal", Value: 9}
	val, _ = exprGtUint.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}
	exprGtFloat := GreaterThanExpression{Field: "FloatVal", Value: 9.5}
	val, _ = exprGtFloat.Evaluate(u)
	if !val {
		t.Errorf("Expected true")
	}

	val, _ = exprGtFloat.Evaluate(nil)
	if val {
		t.Errorf("Expected false")
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
	if err == nil && val != false {
		t.Errorf("Expected error evaluating bad op")
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

	exprGt := GreaterThanExpression{Field: "Name", Value: "alice"}
	val, _ := exprGt.Evaluate(u)
	if !val { t.Errorf("Expected true for gt") }

	exprGte := GreaterThanOrEqualExpression{Field: "Name", Value: "bob"}
	val, _ = exprGte.Evaluate(u)
	if !val { t.Errorf("Expected true for gte") }

	exprLt := LessThanExpression{Field: "Name", Value: "charlie"}
	val, _ = exprLt.Evaluate(u)
	if !val { t.Errorf("Expected true for lt") }

	exprLte := LessThanOrEqualExpression{Field: "Name", Value: "bob"}
	val, _ = exprLte.Evaluate(u)
	if !val { t.Errorf("Expected true for lte") }

	exprGteFalse := GreaterThanOrEqualExpression{Field: "Name", Value: "charlie"}
	val, _ = exprGteFalse.Evaluate(u)
	if val { t.Errorf("Expected false for gte") }

	exprLtFalse := LessThanExpression{Field: "Name", Value: "alice"}
	val, _ = exprLtFalse.Evaluate(u)
	if val { t.Errorf("Expected false for lt") }

	exprLteFalse := LessThanOrEqualExpression{Field: "Name", Value: "alice"}
	val, _ = exprLteFalse.Evaluate(u)
	if val { t.Errorf("Expected false for lte") }

	val, _ = exprGte.Evaluate(nil)
	if val { t.Errorf("Expected false for nil") }
	val, _ = exprLt.Evaluate(nil)
	if val { t.Errorf("Expected false for nil") }
}

func TestComparisonExpressionEvaluateStringsNeq(t *testing.T) {
	u := &testUser{Name: "bob"}

	exprNeq := ComparisonExpression{Operation: "neq", LHS: Field{Name: "Name"}, RHS: Constant{Value: "alice"}}
	val, _ := exprNeq.Evaluate(u)
	if !val { t.Errorf("Expected true for neq") }
}

func TestComparisonExpressionEvaluateStringsCompare(t *testing.T) {
	u := &testUser{Name: "bob"}

	exprGt := ComparisonExpression{Operation: "gt", LHS: Field{Name: "Name"}, RHS: Constant{Value: "alice"}}
	val, _ := exprGt.Evaluate(u)
	if !val { t.Errorf("Expected true for gt string") }

	exprGte := ComparisonExpression{Operation: "gte", LHS: Field{Name: "Name"}, RHS: Constant{Value: "bob"}}
	val, _ = exprGte.Evaluate(u)
	if !val { t.Errorf("Expected true for gte string") }

	exprLt := ComparisonExpression{Operation: "lt", LHS: Field{Name: "Name"}, RHS: Constant{Value: "charlie"}}
	val, _ = exprLt.Evaluate(u)
	if !val { t.Errorf("Expected true for lt string") }

	exprLte := ComparisonExpression{Operation: "lte", LHS: Field{Name: "Name"}, RHS: Constant{Value: "bob"}}
	val, _ = exprLte.Evaluate(u)
	if !val { t.Errorf("Expected true for lte string") }
}

func TestComparisonExpressionEvaluateStringsCompareFailures(t *testing.T) {
	u := &testUser{Name: "bob"}

	exprGt := ComparisonExpression{Operation: "gt", LHS: Field{Name: "Name"}, RHS: Constant{Value: "charlie"}}
	val, _ := exprGt.Evaluate(u)
	if val { t.Errorf("Expected false for gt string") }

	exprGte := ComparisonExpression{Operation: "gte", LHS: Field{Name: "Name"}, RHS: Constant{Value: "charlie"}}
	val, _ = exprGte.Evaluate(u)
	if val { t.Errorf("Expected false for gte string") }

	exprLt := ComparisonExpression{Operation: "lt", LHS: Field{Name: "Name"}, RHS: Constant{Value: "alice"}}
	val, _ = exprLt.Evaluate(u)
	if val { t.Errorf("Expected false for lt string") }

	exprLte := ComparisonExpression{Operation: "lte", LHS: Field{Name: "Name"}, RHS: Constant{Value: "alice"}}
	val, _ = exprLte.Evaluate(u)
	if val { t.Errorf("Expected false for lte string") }
}
