package aliases_test

import (
	"testing"

	aliases "github.com/arran4/go-evaluator/aliases"
)

type user struct {
	Name string
}

func TestAliasesEvaluate(t *testing.T) {
	q := aliases.Q{Expression: &aliases.EQ{Field: "Name", Value: "bob"}}
	if !q.Evaluate(&user{Name: "bob"}) {
		t.Fatalf("expected true")
	}
}
