package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/example/logslice/internal/parser"
	"github.com/example/logslice/internal/runner"
	"github.com/example/logslice/internal/timerange"
)

func main() {
	format := flag.String("format", "nginx", "Log format: nginx, apache, syslog, json, csv")
	start := flag.String("start", "", "Start time (RFC3339)")
	end := flag.String("end", "", "End time (RFC3339)")
	outFmt := flag.String("output", "text", "Output format: text, json")
	jsonField := flag.String("json-time-field", "time", "JSON log timestamp field name")
	csvTimeCol := flag.Int("csv-time-col", 0, "CSV 0-based timestamp column index")
	csvTimeFmt := flag.String("csv-time-fmt", time.RFC3339, "CSV timestamp Go layout")
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
	switch strings.ToLower(*format) {
	case "nginx":
		p = parser.NewNginxParser()
	case "apache":
		p = parser.NewApacheParser()
	case "syslog":
		p = parser.NewSyslogParser()
	case "json":
		p = parser.NewJSONParser(*jsonField)
	case "csv":
		p = parser.NewCSVParser(*csvTimeCol, *csvTimeFmt)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format %q\n", *format)
		os.Exit(1)
	}

	if err := runner.Run(os.Stdin, os.Stdout, p, tr, *outFmt); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
