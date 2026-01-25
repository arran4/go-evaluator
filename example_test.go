package evaluator_test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/go-evaluator/parser/simple"
)

// ExampleQuery_Evaluate demonstrates manual construction of a query and evaluation.
func ExampleQuery_Evaluate() {
	q := evaluator.Query{
		Expression: &evaluator.AndExpression{Expressions: []evaluator.Query{
			{Expression: &evaluator.IsExpression{Field: "Name", Value: "bob"}},
			{Expression: &evaluator.GreaterThanExpression{Field: "Age", Value: 30}},
		}},
	}

	type User struct {
		Name string
		Age  int
	}

	fmt.Println(q.Evaluate(&User{Name: "bob", Age: 35}))
	// Output: true
}

// ExampleQuery_unmarshalJSON shows how to unmarshal a query from JSON.
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

// Example_simpleParser demonstrates parsing a string query.
func Example_simpleParser() {
	queryString := `Category is "Electronics" and Price < 1000`
	query, err := simple.Parse(queryString)
	if err != nil {
		log.Fatal(err)
	}

	type Product struct {
		Name     string
		Category string
		Price    float64
	}

	result := query.Evaluate(&Product{
		Name:     "Laptop",
		Category: "Electronics",
		Price:    999.99,
	})
	fmt.Println(result)
	// Output: true
}

// ExampleFunctionExpression_Evaluate demonstrates how to use the FunctionExpression
// to execute custom logic within the evaluator.
func ExampleFunctionExpression_Evaluate() {
	// 1. Define a custom function
	// This function sums all numeric arguments
	sumFunc := &SumFunction{}

	// 2. Build the expression tree
	// We want to calculate: Sum(10, 20)
	expr := evaluator.FunctionExpression{
		Func: sumFunc,
		Args: []evaluator.Term{
			evaluator.Constant{Value: 10},
			evaluator.Constant{Value: 20},
		},
	}

	// 3. Evaluate the expression
	result, err := expr.Evaluate(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sum: %v\n", result)

	// 4. Nested usage
	// Calculate: Sum(Sum(5, 5), 10) -> Sum(10, 10) -> 20
	nestedExpr := evaluator.FunctionExpression{
		Func: sumFunc,
		Args: []evaluator.Term{
			evaluator.FunctionExpression{
				Func: sumFunc,
				Args: []evaluator.Term{
					evaluator.Constant{Value: 5},
					evaluator.Constant{Value: 5},
				},
			},
			evaluator.Constant{Value: 10},
		},
	}

	resultNested, err := nestedExpr.Evaluate(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Nested Sum: %v\n", resultNested)

	// Output:
	// Sum: 30
	// Nested Sum: 20
}

// SumFunction implements evaluator.Function
type SumFunction struct{}

func (s *SumFunction) Call(args ...interface{}) (interface{}, error) {
	sum := 0.0
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			sum += float64(v)
		case float64:
			sum += v
		default:
			return nil, fmt.Errorf("unsupported type %T", arg)
		}
	}
	return sum, nil
}
