package simple

import (
	"reflect"
	"testing"
)

type testUser struct {
	Name  string
	Age   int
	Tags  []string
	Score float64
}

func TestParseAndEvaluate(t *testing.T) {
	expr := `Name is "bob" and Age > 30`
	q, err := Parse(expr)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	u := &testUser{Name: "bob", Age: 35}
	if v, err := q.Evaluate(u); err != nil || !v {
		t.Errorf("evaluation failed: %v %v", v, err)
	}
}

func TestRoundTrip(t *testing.T) {
	exprs := []string{
		`Name is "bob"`,
		`Name is not "alice"`,
		`Score >= 4.5`,
		`Tags contains "go"`,
		`not (Name is "alice")`,
		`(Name is "bob" and Age > 30) or Score < 2`,
	}
	for _, e := range exprs {
		q, err := Parse(e)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		s := Stringify(q)
		q2, err := Parse(s)
		if err != nil {
			t.Fatalf("parse round: %v", err)
		}
		if !reflect.DeepEqual(q, q2) {
			t.Errorf("round trip mismatch for %s", e)
		}
	}
}
