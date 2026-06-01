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

func TestParserErrors(t *testing.T) {
	cases := []string{
		`Name is`,
		`Name is "bob" and`,
		`not`,
		`(Name is "bob"`,
		`"bob" is Name`, // this might work or fail
		``,
		`(Name is "bob"`, // missing )
		`Name > `, // missing right side
	}
	for _, c := range cases {
		_, err := Parse(c)
		if err == nil {
			t.Errorf("Expected error for %q", c)
		}
	}
}

func TestStringify(t *testing.T) {
	exprs := []string{
		`Name >= 4`,
		`Name > 4`,
		`Name < 4`,
		`Name <= 4`,
		`Name contains "bob"`,
		`Name is "bob"`,
		`Name is not "bob"`,
		`not (Name is "bob")`,
		`(Name > 4 and Name < 10)`,
		`(Name > 4 or Name < 10)`,
	}
	for _, e := range exprs {
		q, err := Parse(e)
		if err != nil {
			t.Errorf("Parse error for %q: %v", e, err)
			continue
		}
		s := Stringify(q)
		if s == "" {
			t.Errorf("Stringify empty for %q", e)
		}
	}
}

func TestValToString(t *testing.T) {
	cases := []struct{
		val interface{}
		expect string
	}{
		{"bob", `"bob"`},
		{4, "4"},
		{4.5, "4.5"},
		{true, "true"},
		{[]int{1}, "[1]"},
	}
	for _, c := range cases {
		s := valToString(c.val)
		if s != c.expect {
			t.Errorf("Expected %q, got %q", c.expect, s)
		}
	}
}
