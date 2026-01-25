package evaluator

import (
	"encoding/json"
	"testing"
)

type testUser struct {
	Name  string
	Age   int
	Tags  []string
	Score float64
}

func TestContainsExpression(t *testing.T) {
	u := &testUser{Tags: []string{"a", "b"}}
	if v, err := (ContainsExpression{Field: "Tags", Value: "a"}.Evaluate(u)); err != nil || !v {
		t.Errorf("expected true, got %v, %v", v, err)
	}
	if v, err := (ContainsExpression{Field: "Tags", Value: "c"}.Evaluate(u)); err != nil || v {
		t.Errorf("expected false, got %v, %v", v, err)
	}
}

func TestIsAndIsNot(t *testing.T) {
	u := &testUser{Name: "bob"}
	if v, err := (IsExpression{Field: "Name", Value: "bob"}.Evaluate(u)); err != nil || !v {
		t.Errorf("is failed: %v %v", v, err)
	}
	if v, err := (IsNotExpression{Field: "Name", Value: "alice"}.Evaluate(u)); err != nil || !v {
		t.Errorf("isnot failed: %v %v", v, err)
	}
}

func TestComparisons(t *testing.T) {
	u := &testUser{Age: 40, Score: 4.5}
	if v, err := (&GreaterThanExpression{Field: "Age", Value: 30}).Evaluate(u); err != nil || !v {
		t.Errorf("gt failed: %v %v", v, err)
	}
	if v, err := (&GreaterThanOrEqualExpression{Field: "Age", Value: 40}).Evaluate(u); err != nil || !v {
		t.Errorf("gte failed: %v %v", v, err)
	}
	if v, err := (&LessThanExpression{Field: "Score", Value: 5}).Evaluate(u); err != nil || !v {
		t.Errorf("lt failed: %v %v", v, err)
	}
	if v, err := (&LessThanOrEqualExpression{Field: "Score", Value: 4.5}).Evaluate(u); err != nil || !v {
		t.Errorf("lte failed: %v %v", v, err)
	}

	if v, err := (&GreaterThanExpression{Field: "Missing", Value: 1}).Evaluate(u); err != nil || v {
		t.Errorf("gt missing field should be false: %v %v", v, err)
	}
}

func TestStringComparisons(t *testing.T) {
	u := &testUser{Name: "bob"}

	if v, err := (&GreaterThanExpression{Field: "Name", Value: "ann"}).Evaluate(u); err != nil || !v {
		t.Errorf("gt string failed: %v %v", v, err)
	}
	if v, err := (&GreaterThanOrEqualExpression{Field: "Name", Value: "bob"}).Evaluate(u); err != nil || !v {
		t.Errorf("gte string failed: %v %v", v, err)
	}
	if v, err := (&LessThanExpression{Field: "Name", Value: "carol"}).Evaluate(u); err != nil || !v {
		t.Errorf("lt string failed: %v %v", v, err)
	}
	if v, err := (&LessThanOrEqualExpression{Field: "Name", Value: "bob"}).Evaluate(u); err != nil || !v {
		t.Errorf("lte string failed: %v %v", v, err)
	}
}

func TestLogicalExpressions(t *testing.T) {
	u := &testUser{Name: "bob", Age: 41}
	and := AndExpression{Expressions: []Query{
		{Expression: &IsExpression{Field: "Name", Value: "bob"}},
		{Expression: &GreaterThanExpression{Field: "Age", Value: 40}},
	}}
	if v, err := and.Evaluate(u); err != nil || !v {
		t.Errorf("and failed: %v %v", v, err)
	}
	or := OrExpression{Expressions: []Query{
		{Expression: &IsExpression{Field: "Name", Value: "alice"}},
		{Expression: &GreaterThanExpression{Field: "Age", Value: 40}},
	}}
	if v, err := or.Evaluate(u); err != nil || !v {
		t.Errorf("or failed: %v %v", v, err)
	}
	not := NotExpression{Expression: Query{Expression: &IsExpression{Field: "Name", Value: "alice"}}}
	if v, err := not.Evaluate(u); err != nil || !v {
		t.Errorf("not failed: %v %v", v, err)
	}
}

func TestNonPointerInput(t *testing.T) {
	u := testUser{Tags: []string{"a"}, Name: "bob"}
	if v, err := (ContainsExpression{Field: "Tags", Value: "a"}).Evaluate(u); err != nil || v {
		t.Errorf("expected false for non-pointer input: %v %v", v, err)
	}
	if v, err := (IsExpression{Field: "Name", Value: "bob"}).Evaluate(u); err != nil || v {
		t.Errorf("expected false for non-pointer input: %v %v", v, err)
	}
}

func TestQueryUnmarshalAndEvaluate(t *testing.T) {
	js := `{
        "Expression": {
            "Type": "And",
            "Expression": {
                "Expressions": [
                    {"Expression": {"Type": "Is", "Expression": {"Field": "Name", "Value": "bob"}}},
                    {"Expression": {"Type": "GT", "Expression": {"Field": "Age", "Value": 30}}}
                ]
            }
        }
    }`
	var q Query
	if err := json.Unmarshal([]byte(js), &q); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	u := &testUser{Name: "bob", Age: 35}
	if v, err := q.Evaluate(u); err != nil || !v {
		t.Errorf("query evaluate failed: %v %v", v, err)
	}
}

func TestQueryUnmarshalStringCompare(t *testing.T) {
	js := `{
        "Expression": {
            "Type": "LT",
            "Expression": {
                "Field": "Name",
                "Value": "d"
            }
        }
    }`
	var q Query
	if err := json.Unmarshal([]byte(js), &q); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	u := &testUser{Name: "bob"}
	if v, err := q.Evaluate(u); err != nil || !v {
		t.Errorf("string comparison in query failed: %v %v", v, err)
	}
}
func TestQueryMarshalRoundTrip(t *testing.T) {
	q := Query{Expression: &AndExpression{Expressions: []Query{
		{Expression: &IsExpression{Field: "Name", Value: "bob"}},
		{Expression: &GreaterThanExpression{Field: "Age", Value: 30}},
	}}}
	b1, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var q2 Query
	if err := json.Unmarshal(b1, &q2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	b2, err := json.Marshal(q2)
	if err != nil {
		t.Fatalf("marshal2: %v", err)
	}
	if string(b1) != string(b2) {
		t.Errorf("round trip json mismatch\norig: %s\nback: %s", b1, b2)
	}
}

func TestQueryMarshalEvaluate(t *testing.T) {
	q := Query{Expression: &NotExpression{Expression: Query{Expression: &IsExpression{Field: "Name", Value: "alice"}}}}
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var q2 Query
	if err := json.Unmarshal(b, &q2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	u := &testUser{Name: "bob"}
	v1, err1 := q.Evaluate(u)
	v2, err2 := q2.Evaluate(u)
	if err1 != err2 {
		t.Fatalf("error mismatch: %v vs %v", err1, err2)
	}
	if v1 != v2 {
		t.Errorf("evaluation mismatch after round trip: %v vs %v", v1, v2)
	}
}
