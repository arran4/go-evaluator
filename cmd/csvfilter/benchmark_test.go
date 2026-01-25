package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"testing"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/go-evaluator/parser/simple"
)

type trueExpression struct{}

func (e trueExpression) Evaluate(i interface{}) bool {
	return true
}

func BenchmarkProcess(b *testing.B) {
	// Setup input
	headers := []string{"id", "name", "age", "city", "active"}
	var buf bytes.Buffer
	// write headers
	for i, h := range headers {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(h)
	}
	buf.WriteString("\n")

	// write rows
	row := "1,John Doe,30,New York,true\n"
	for i := 0; i < 1000; i++ {
		buf.WriteString(row)
	}
	inputData := buf.Bytes()

	// Mock query
	q := evaluator.Query{
		Expression: trueExpression{},
	}

	// Redirect stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	null, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Fatal(err)
	}
	os.Stdout = null

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(inputData)
		wh := true
		if err := process(r, q, &wh); err != nil {
			b.Fatal(err)
		}
	}
}

func TestProcess_Functional(t *testing.T) {
	input := `id,name,age
1,Alice,30
2,Bob,25
3,Charlie,35
`
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Query: age > 28
	q, err := simple.Parse("age > 28")
	if err != nil {
		t.Fatalf("Failed to parse query: %v", err)
	}

	reader := bytes.NewReader([]byte(input))
	wh := true

	errChan := make(chan error, 1)
	go func() {
		errChan <- process(reader, q, &wh)
		w.Close()
	}()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	if err := <-errChan; err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	output := buf.String()

	// Verify output
	cr := csv.NewReader(bytes.NewReader(buf.Bytes()))
	rows, err := cr.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse output CSV: %v", err)
	}

	// Expect header + Alice (30) + Charlie (35)
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows (header + 2 matches), got %d. Output:\n%s", len(rows), output)
	}
	if len(rows) > 1 && rows[1][1] != "Alice" {
		t.Errorf("Expected first match to be Alice, got %s", rows[1][1])
	}
	if len(rows) > 2 && rows[2][1] != "Charlie" {
		t.Errorf("Expected second match to be Charlie, got %s", rows[2][1])
	}
}
