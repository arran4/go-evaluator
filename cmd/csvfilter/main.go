package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/arran4/go-evaluator"
	"github.com/arran4/go-evaluator/parser/simple"
)

func process(r io.Reader, q evaluator.Query, writeHeader *bool) error {
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
	m := make(map[string]interface{}, len(headers))
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		clear(m)
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -e <expression> [file ...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Filter CSV rows matching the expression. If no files are given, input is read from standard input.")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	expr := flag.String("e", "", "expression to apply to each row")
	flag.Parse()
	if *expr == "" {
		log.Fatal("-e expression required")
	}
	q, err := simple.Parse(*expr)
	if err != nil {
		log.Fatalf("parse expression: %v", err)
	}
	files := flag.Args()
	writeHeader := true
	if len(files) == 0 {
		if err := process(os.Stdin, q, &writeHeader); err != nil {
			log.Fatal(err)
		}
		return
	}
	for _, f := range files {
		fh, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		if err := process(fh, q, &writeHeader); err != nil {
			_ = fh.Close()
			log.Fatal(err)
		}
		_ = fh.Close()
	}
}
