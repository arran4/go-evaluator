package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/go-evaluator/parser/simple"
)

func process(r io.Reader, q evaluator.Query) error {
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -e <expression> [file ...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Filter JSON Lines records matching the expression. Reads from standard input when no files are provided.")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	expr := flag.String("e", "", "expression to apply to each object")
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
		if err := process(os.Stdin, q); err != nil {
			log.Fatal(err)
		}
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		if err := process(fh, q); err != nil {
			fh.Close()
			log.Fatal(err)
		}
		fh.Close()
	}
}
