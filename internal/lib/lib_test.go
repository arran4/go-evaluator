package lib

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/arran4/go-evaluator/parser/simple"
)

func TestProcessCSV(t *testing.T) {
	input := `name,age
alice,30
bob,25
charlie,35`
	expr := "age > 28"

	q, err := simple.Parse(expr)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	r := bytes.NewBufferString(input)
	var w bytes.Buffer
	writeHeader := true

	if err := processCSV(r, &w, q, &writeHeader); err != nil {
		t.Fatalf("processCSV error: %v", err)
	}

	expected := "name,age\nalice,30\ncharlie,35\n"
	if w.String() != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, w.String())
	}
}

func BenchmarkProcessCSV(b *testing.B) {
	// Prepare a large-ish CSV input
	var buf bytes.Buffer
	buf.WriteString("id,value,category\n")
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(&buf, "%d,%d,cat%d\n", i, i%100, i%3)
	}
	inputData := buf.Bytes()

	expr := "value > 50"
	q, err := simple.Parse(expr)
	if err != nil {
		b.Fatalf("parse error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(inputData)
		writeHeader := true
		if err := processCSV(r, io.Discard, q, &writeHeader); err != nil {
			b.Fatalf("processCSV error: %v", err)
		}
	}
}

func TestEvaluateJSON(t *testing.T) {
	input := `{"name": "alice", "age": 30}`
	q, _ := simple.Parse(`name is "alice"`)
	r := bytes.NewReader([]byte(input))
	val, err := evaluateJSON(r, q)
	if err != nil {
		t.Fatalf("evaluateJSON error: %v", err)
	}
	if !val {
		t.Errorf("Expected true")
	}
}

func TestEvaluateYAML(t *testing.T) {
	input := `name: alice
age: 30`
	q, _ := simple.Parse(`name is "alice"`)
	r := bytes.NewReader([]byte(input))
	val, err := evaluateYAML(r, q)
	if err != nil {
		t.Fatalf("evaluateYAML error: %v", err)
	}
	if !val {
		t.Errorf("Expected true")
	}
}

func TestProcessJSONL(t *testing.T) {
	input := `{"name": "alice", "age": 30}
{"name": "bob", "age": 25}`
	q, _ := simple.Parse(`age > 28`)
	r := bytes.NewReader([]byte(input))
	err := processJSONL(r, q)
	if err != nil {
		t.Fatalf("processJSONL error: %v", err)
	}
}

func TestProcessJSONLEOF(t *testing.T) {
	input := `{"name": "alice", "age": 30}
`
	q, _ := simple.Parse(`age > 28`)
	r := bytes.NewReader([]byte(input))
	err := processJSONL(r, q)
	if err != nil {
		t.Fatalf("processJSONL error: %v", err)
	}
}
