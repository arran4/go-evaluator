package lib

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/go-evaluator/parser/simple"
)

// CsvFilter filters CSV rows matching the expression.
func CsvFilter(expr string, files ...string) {
	if expr == "" {
		log.Fatal("-e expression required")
	}
	q, err := simple.Parse(expr)
	if err != nil {
		log.Fatalf("parse expression: %v", err)
	}
	writeHeader := true
	if len(files) == 0 {
		if err := processCSV(os.Stdin, q, &writeHeader); err != nil {
			log.Fatal(err)
		}
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		if err := processCSV(fh, q, &writeHeader); err != nil {
			fh.Close()
			log.Fatal(err)
		}
		fh.Close()
	}
}

func processCSV(r io.Reader, q evaluator.Query, writeHeader *bool) error {
	cr := csv.NewReader(r)
	headers, err := cr.Read()
	if err != nil {
		return err
	}
	cw := csv.NewWriter(os.Stdout)
	if *writeHeader {
		if err := cw.Write(headers); err != nil {
			return err
		}
		*writeHeader = false
	}
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		m := make(map[string]interface{}, len(headers))
		for i, h := range headers {
			if i < len(rec) {
				m[h] = rec[i]
			}
		}
		if q.Evaluate(m) {
			if err := cw.Write(rec); err != nil {
				return err
			}
		}
	}
	cw.Flush()
	return cw.Error()
}

// JsonlFilter filters JSON Lines records matching the expression.
func JsonlFilter(expr string, files ...string) {
	if expr == "" {
		log.Fatal("-e expression required")
	}
	q, err := simple.Parse(expr)
	if err != nil {
		log.Fatalf("parse expression: %v", err)
	}
	if len(files) == 0 {
		if err := processJSONL(os.Stdin, q); err != nil {
			log.Fatal(err)
		}
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		if err := processJSONL(fh, q); err != nil {
			fh.Close()
			log.Fatal(err)
		}
		fh.Close()
	}
}

func processJSONL(r io.Reader, q evaluator.Query) error {
	dec := json.NewDecoder(r)
	enc := json.NewEncoder(os.Stdout)
	for {
		var m map[string]interface{}
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if q.Evaluate(m) {
			if err := enc.Encode(m); err != nil {
				return err
			}
		}
	}
	return nil
}

// JsonTest evaluates a JSON document against the expression.
func JsonTest(expr string, files ...string) {
	if expr == "" {
		log.Fatal("-e expression required")
	}
	q, err := simple.Parse(expr)
	if err != nil {
		log.Fatalf("parse expression: %v", err)
	}
	if len(files) == 0 {
		ok, err := evaluateJSON(os.Stdin, q)
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			return
		}
		os.Exit(1)
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		ok, err := evaluateJSON(fh, q)
		fh.Close()
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			os.Exit(1)
		}
	}
}

func evaluateJSON(r io.Reader, q evaluator.Query) (bool, error) {
	dec := json.NewDecoder(r)
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return false, err
	}
	return q.Evaluate(m), nil
}

// YamlTest evaluates a YAML document against the expression.
func YamlTest(expr string, files ...string) {
	if expr == "" {
		log.Fatal("-e expression required")
	}
	q, err := simple.Parse(expr)
	if err != nil {
		log.Fatalf("parse expression: %v", err)
	}
	if len(files) == 0 {
		ok, err := evaluateYAML(os.Stdin, q)
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			return
		}
		os.Exit(1)
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		ok, err := evaluateYAML(fh, q)
		fh.Close()
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			os.Exit(1)
		}
	}
}

func evaluateYAML(r io.Reader, q evaluator.Query) (bool, error) {
	dec := yaml.NewDecoder(r)
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return false, err
	}
	return q.Evaluate(m), nil
}
