package main

import (
	"github.com/arran4/go-evaluator/internal/lib"
)

// CsvFilter is a subcommand `evaluator csvfilter`
// Flags:
//
//	expr: -e Expression
//	files: ... Files
func CsvFilter(expr string, files ...string) {
	lib.CsvFilter(expr, files...)
}

// JsonlFilter is a subcommand `evaluator jsonlfilter`
// Flags:
//
//	expr: -e Expression
//	files: ... Files
func JsonlFilter(expr string, files ...string) {
	lib.JsonlFilter(expr, files...)
}

// JSONTest is a subcommand `evaluator jsontest`
// Flags:
//
//	expr: -e Expression
//	files: ... Files
func JSONTest(expr string, files ...string) {
	lib.JSONTest(expr, files...)
}

// YamlTest is a subcommand `evaluator yamltest`
// Flags:
//
//	expr: -e Expression
//	files: ... Files
func YamlTest(expr string, files ...string) {
	lib.YamlTest(expr, files...)
}

//go:generate go run github.com/arran4/go-subcommand/cmd/gosubc generate --dir ../..
