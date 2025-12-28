# Evaluator

The **evaluator** package implements a small expression language for querying Go
structs. Expressions are represented as Go structs that can be combined using
logical operators. Comparison expressions support both numeric and string
values.

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

## Command-line Tools

The project includes small utilities for working with common data formats.

- **csvfilter** – filters CSV rows. Usage:
  `csvfilter -e '<expression>' [file ...]`
  When no files are specified, input is read from standard input.
- **jsonlfilter** – filters newline-delimited JSON records. It accepts the same
  arguments as `csvfilter` and writes matching JSON objects to standard output.
- **jsontest** – evaluates a single JSON document. It exits with status 0 if the
  expression matches and 1 otherwise. Multiple files can be supplied and all
  must satisfy the expression. With no files it reads from standard input.
- **yamltest** – like `jsontest` but for YAML documents.

Expressions use the simple syntax implemented by the parser in
`parser/simple`. For example:

```bash
jsonlfilter -e 'GT(age,30) and Is(country,"US")' users.jsonl
```

## Running Tests

Run `go test ./...` to execute the unit tests.


## License

This project is licensed under the [MIT License](LICENSE).

