# Evaluator

The **evaluator** package implements a small expression language for querying Go
structs. Expressions are represented as Go structs that can be combined using
logical operators. Comparison expressions support both numeric and string
values.

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

## Running Tests

Run `go test ./...` to execute the unit tests.

## Extracting the Package

A helper script `extract_evaluator.sh` is provided in the repository root. It
creates a new repository containing only the evaluator history.

```sh
./extract_evaluator.sh ../evaluator-repo
```

