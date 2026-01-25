package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/go-evaluator/parser/simple"
)

func evaluate(r io.Reader, q evaluator.Query) (bool, error) {
	dec := yaml.NewDecoder(r)
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		return false, err
	}
	return q.Evaluate(m)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -e <expression> [file ...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Evaluate a YAML document against the expression. Reads from stdin when no files are specified. Exits with status 1 if the expression does not match.")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	expr := flag.String("e", "", "expression to test against the document")
	flag.Parse()
	if *expr == "" {
		log.Fatal("-e expression required")
	}
	q, err := simple.Parse(*expr)
	if err != nil {
		log.Fatalf("parse expression: %v", err)
	}
	files := flag.Args()
	if len(files) == 0 {
		ok, err := evaluate(os.Stdin, q)
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
		ok, err := evaluate(fh, q)
		_ = fh.Close()
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			os.Exit(1)
		}
	}
}
