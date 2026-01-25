package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/arran4/go-evaluator/parser/simple"
)

func BenchmarkProcess(b *testing.B) {
	// Prepare input data
	var buf bytes.Buffer
	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			buf.WriteString(`{"Name": "match", "Age": 30, "Extra": "data", "Nested": {"Key": "Value"}}` + "\n")
		} else {
			buf.WriteString(`{"Name": "other", "Age": 30, "Extra": "data", "Nested": {"Key": "Value"}}` + "\n")
		}
	}
	input := buf.Bytes()

	q, err := simple.Parse(`Name is "match"`)
	if err != nil {
		b.Fatalf("parse error: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		err := process(r, io.Discard, q)
		if err != nil {
			b.Fatalf("process error: %v", err)
		}
	}
}

func TestProcess(t *testing.T) {
	input := `{"Name": "match"}
{"Name": "other"}
{"Name": "match"}
`
	expected := `{"Name":"match"}
{"Name":"match"}
`
	q, err := simple.Parse(`Name is "match"`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var out bytes.Buffer
	err = process(bytes.NewBufferString(input), &out, q)
	if err != nil {
		t.Fatalf("process error: %v", err)
	}

	if out.String() != expected {
		t.Errorf("expected output:\n%s\ngot:\n%s", expected, out.String())
	}
}
