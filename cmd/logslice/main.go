package main

import (
	"flag"
	"fmt"
	"os"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/parser"
	"logslice/internal/timerange"
)

func main() {
	format := flag.String("format", "nginx", "Log format: nginx or apache")
	start := flag.String("start", "", "Start time (RFC3339, e.g. 2024-01-01T00:00:00Z)")
	end := flag.String("end", "", "End time (RFC3339, e.g. 2024-01-02T00:00:00Z)")
	outputFmt := flag.String("output", "text", "Output format: text or json")
	file := flag.String("file", "", "Log file to parse (default: stdin)")
	flag.Parse()

	if *start == "" || *end == "" {
		fmt.Fprintln(os.Stderr, "error: --start and --end are required")
		flag.Usage()
		os.Exit(1)
	}

	tr, err := timerange.New(*start, *end)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid time range: %v\n", err)
		os.Exit(1)
	}

	var p parser.Parser
	switch *format {
	case "nginx":
		p = parser.NewNginxParser()
	case "apache":
		p = parser.NewApacheParser()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format %q\n", *format)
		os.Exit(1)
	}

	var in *os.File
	if *file == "" {
		in = os.Stdin
	} else {
		in, err = os.Open(*file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		defer in.Close()
	}

	f := filter.New(tr)
	fmt := output.NewFormatter(*outputFmt, os.Stdout)

	if err := run(in, p, f, fmt); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
