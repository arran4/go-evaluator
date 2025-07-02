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
	if !(ContainsExpression{Field: "Tags", Value: "a"}.Evaluate(u)) {
		t.Errorf("expected true")
	}
	if (ContainsExpression{Field: "Tags", Value: "c"}.Evaluate(u)) {
		t.Errorf("expected false")
	}
}

func TestIsAndIsNot(t *testing.T) {
	u := &testUser{Name: "bob"}
	if !(IsExpression{Field: "Name", Value: "bob"}.Evaluate(u)) {
		t.Errorf("is failed")
	}
	if !(IsNotExpression{Field: "Name", Value: "alice"}.Evaluate(u)) {
		t.Errorf("isnot failed")
	}
}

func TestComparisons(t *testing.T) {
	u := &testUser{Age: 40, Score: 4.5}
	if !(GreaterThanExpression{Field: "Age", Value: 30}.Evaluate(u)) {
		t.Errorf("gt failed")
	}
	if !(GreaterThanOrEqualExpression{Field: "Age", Value: 40}.Evaluate(u)) {
		t.Errorf("gte failed")
	}
	if !(LessThanExpression{Field: "Score", Value: 5}.Evaluate(u)) {
		t.Errorf("lt failed")
	}
	if !(LessThanOrEqualExpression{Field: "Score", Value: 4.5}.Evaluate(u)) {
		t.Errorf("lte failed")
	}

	if (GreaterThanExpression{Field: "Missing", Value: 1}.Evaluate(u)) {
		t.Errorf("gt missing field should be false")
	}
}

func TestStringComparisons(t *testing.T) {
	u := &testUser{Name: "bob"}

	if !(GreaterThanExpression{Field: "Name", Value: "ann"}.Evaluate(u)) {
		t.Errorf("gt string failed")
	}
	if !(GreaterThanOrEqualExpression{Field: "Name", Value: "bob"}.Evaluate(u)) {
		t.Errorf("gte string failed")
	}
	if !(LessThanExpression{Field: "Name", Value: "carol"}.Evaluate(u)) {
		t.Errorf("lt string failed")
	}
	if !(LessThanOrEqualExpression{Field: "Name", Value: "bob"}.Evaluate(u)) {
		t.Errorf("lte string failed")
	}
}

func TestLogicalExpressions(t *testing.T) {
	u := &testUser{Name: "bob", Age: 41}
	and := AndExpression{Expressions: []Query{
		{Expression: &IsExpression{Field: "Name", Value: "bob"}},
		{Expression: &GreaterThanExpression{Field: "Age", Value: 40}},
	}}
	if !(and.Evaluate(u)) {
		t.Errorf("and failed")
	}
	or := OrExpression{Expressions: []Query{
		{Expression: &IsExpression{Field: "Name", Value: "alice"}},
		{Expression: &GreaterThanExpression{Field: "Age", Value: 40}},
	}}
	if !(or.Evaluate(u)) {
		t.Errorf("or failed")
	}
	not := NotExpression{Expression: Query{Expression: &IsExpression{Field: "Name", Value: "alice"}}}
	if !(not.Evaluate(u)) {
		t.Errorf("not failed")
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
	if !(q.Evaluate(u)) {
		t.Errorf("query evaluate failed")
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
	if !q.Evaluate(u) {
		t.Errorf("string comparison in query failed")
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
	if q.Evaluate(u) != q2.Evaluate(u) {
		t.Errorf("evaluation mismatch after round trip")
	}
}
