package evaluator_test

import (
	"encoding/json"
	"fmt"

	"github.com/arran4/go-evaluator"
)

// Example demonstrating manual construction of a query and evaluation.
func ExampleQuery_Evaluate() {
	type User struct {
		Name string
		Age  int
	}
	q := evaluator.Query{
		Expression: &evaluator.AndExpression{Expressions: []evaluator.Query{
			{Expression: &evaluator.IsExpression{Field: "Name", Value: "bob"}},
			{Expression: &evaluator.GreaterThanExpression{Field: "Age", Value: 30}},
		}},
	}
	fmt.Println(q.Evaluate(&User{Name: "bob", Age: 35}))
	// Output: true
}

// Example showing how to unmarshal a query from JSON.
func ExampleQuery_unmarshalJSON() {
	js := `{
        "Expression": {
            "Type": "Contains",
            "Expression": {
                "Field": "Tags",
                "Value": "go"
            }
        }
    }`
	var q evaluator.Query
	_ = json.Unmarshal([]byte(js), &q)
	type Post struct{ Tags []string }
	fmt.Println(q.Evaluate(&Post{Tags: []string{"go", "news"}}))
	// Output: true
}
