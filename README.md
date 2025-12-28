# Evaluator

The **evaluator** package implements a small expression language for querying Go
structs. Expressions are represented as Go structs that can be combined using
logical operators. Comparison expressions support both numeric and string
values.

## Why use this?

The Evaluator library allows you to:
- **Dynamic Filtering**: Define filtering logic at runtime (e.g., from configuration files or user input) rather than hardcoding it.
- **Safe Querying**: Expose a simple, safe query capability to end-users without exposing full SQL or code execution.
- **Portability**: Serialize queries to JSON to store them in a database or send them over a network.
- **Type Safety**: Works with standard Go structs and types.

## Installation

To use the library in your Go project:

```bash
go get github.com/arran4/go-evaluator
```

To install the command-line tools:

```bash
go install github.com/arran4/go-evaluator/cmd/csvfilter@latest
go install github.com/arran4/go-evaluator/cmd/jsonlfilter@latest
go install github.com/arran4/go-evaluator/cmd/jsontest@latest
go install github.com/arran4/go-evaluator/cmd/yamltest@latest
```

## Features

- Equality and inequality checks (`Is`, `IsNot`)
- Numeric and lexical comparisons (`GT`, `GTE`, `LT`, `LTE`)
- Membership checks with `Contains`
- Logical composition using `And`, `Or` and `Not`
- JSON serialisation for easy storage or transmission of queries

## Basic Usage

Create a query using the provided expression types and call `Evaluate` with your
target struct:

```go
q := evaluator.Query{
    Expression: &evaluator.AndExpression{Expressions: []evaluator.Query{
        {Expression: &evaluator.IsExpression{Field: "Name", Value: "bob"}},
        {Expression: &evaluator.GreaterThanExpression{Field: "Age", Value: 30}},
    }},
}

matched := q.Evaluate(&User{Name: "bob", Age: 35})
```

## Integration Example

This example demonstrates how to integrate the evaluator into an application to
filter a list of structs based on a dynamic query string (e.g., from user input).

```go
package main

import (
	"fmt"
	"log"

	"github.com/arran4/go-evaluator/parser/simple"
)

type Product struct {
	Name     string
	Category string
	Price    float64
	InStock  bool
}

func main() {
	// 1. Data source
	products := []Product{
		{"Laptop", "Electronics", 999.99, true},
		{"Coffee Mug", "Kitchen", 12.50, true},
		{"Headphones", "Electronics", 49.99, false},
	}

	// 2. Query (could come from user input, API, config, etc.)
	// Find all Electronics under $1000
	queryString := `Category is "Electronics" and Price < 1000`

	// 3. Parse the query
	query, err := simple.Parse(queryString)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Filter the list
	var filtered []Product
	for _, p := range products {
		if query.Evaluate(&p) {
			filtered = append(filtered, p)
		}
	}

	// 5. Use results
	for _, p := range filtered {
		fmt.Printf("Found: %s ($%.2f)\n", p.Name, p.Price)
	}
}
```

## JSON Queries

Queries can be marshalled to and from JSON. This is handy for configuration
files or network APIs.

```go
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
if err := json.Unmarshal([]byte(js), &q); err != nil {
    log.Fatal(err)
}
```

## Expression Guide

Each query expression implements the `Expression` interface. The table below
lists the available types and their purpose:

| Type                    | Purpose                                         |
|-------------------------|-------------------------------------------------|
| `Is` / `IsNot`          | Check equality or inequality of a field         |
| `GT` / `GTE`            | Numeric or lexical "greater than" comparisons   |
| `LT` / `LTE`            | Numeric or lexical "less than" comparisons      |
| `Contains`              | Test that a slice field contains a value        |
| `And` / `Or` / `Not`    | Compose other expressions logically             |

Example usage:

```go
evaluator.Query{Expression: &evaluator.NotExpression{Expression: evaluator.Query{
    Expression: &evaluator.IsExpression{Field: "Deleted", Value: true},
}}}
```

## CLI Usage & Syntax

The command-line tools use a simple string syntax to define expressions.

**Operators:**
- `is`, `is not`: Equality checks
- `>`, `>=`, `<`, `<=`: Numeric/Lexical comparison
- `contains`: Checks if a list contains a value
- `and`, `or`, `not`: Logical operators
- `(...)`: Grouping

**Values:**
- Strings: `"value"`
- Numbers: `123`, `45.67`
- Booleans: `true`, `false`

**Examples:**
- `Status is "active"`
- `Age >= 18`
- `Tags contains "admin"`
- `(Role is "admin" or Role is "moderator") and Active is true`

## Command-line Tools

The project includes small utilities for working with common data formats.

### csvfilter
Filters CSV rows based on headers.

**Usage:**
```bash
# Given data.csv:
# name,age,city
# alice,30,ny
# bob,25,sf

csvfilter -e 'age > 28' data.csv
# Output:
# name,age,city
# alice,30,ny
```

### jsonlfilter
Filters newline-delimited JSON records.

**Usage:**
```bash
# Given logs.jsonl:
# {"level":"info", "msg":"started"}
# {"level":"error", "msg":"failed"}

jsonlfilter -e 'level is "error"' logs.jsonl
# Output:
# {"level":"error", "msg":"failed"}
```

### jsontest
Evaluates a single JSON document (or multiple files). Returns exit code 0 on match, 1 otherwise.

**Usage:**
```bash
# Check if config.json is valid for production
jsontest -e 'environment is "production" and debug is false' config.json
if [ $? -eq 0 ]; then
    echo "Production ready"
fi
```

### yamltest
Like `jsontest` but for YAML documents.

**Usage:**
```bash
yamltest -e 'replicas >= 3' deployment.yaml
```

## Running Tests

Run `go test ./...` to execute the unit tests.

## License

This project is licensed under the [MIT License](LICENSE).
