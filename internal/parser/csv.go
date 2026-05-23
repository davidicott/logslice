package parser

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

// csvEntry represents a single parsed CSV log line.
type csvEntry struct {
	ts     time.Time
	fields []string
	headers []string
}

func (e *csvEntry) Timestamp() time.Time { return e.ts }

func (e *csvEntry) String() string {
	pairs := make([]string, 0, len(e.headers))
	for i, h := range e.headers {
		if i < len(e.fields) {
			pairs = append(pairs, fmt.Sprintf("%s=%s", h, e.fields[i]))
		}
	}
	return strings.Join(pairs, " ")
}

// CSVParser parses CSV log files with a configurable timestamp column.
type CSVParser struct {
	timeColumn  int
	timeFormat  string
	headers     []string
	headerParsed bool
}

// NewCSVParser creates a CSVParser. timeColumn is the 0-based index of the
// timestamp field; timeFormat is the Go time layout string.
func NewCSVParser(timeColumn int, timeFormat string) *CSVParser {
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	return &CSVParser{
		timeColumn: timeColumn,
		timeFormat: timeFormat,
	}
}

// Parse implements Parser. The first non-empty line is treated as a header row.
func (p *CSVParser) Parse(line string) (Entry, error) {
	if strings.TrimSpace(line) == "" {
		return nil, nil
	}

	r := csv.NewReader(strings.NewReader(line))
	fields, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("csv parse error: %w", err)
	}

	if !p.headerParsed {
		p.headers = fields
		p.headerParsed = true
		return nil, nil
	}

	if p.timeColumn >= len(fields) {
		return nil, fmt.Errorf("time column %d out of range (got %d fields)", p.timeColumn, len(fields))
	}

	ts, err := time.Parse(p.timeFormat, strings.TrimSpace(fields[p.timeColumn]))
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %q: %w", fields[p.timeColumn], err)
	}

	headers := p.headers
	if len(headers) == 0 {
		headers = make([]string, len(fields))
		for i := range fields {
			headers[i] = fmt.Sprintf("col%d", i)
		}
	}

	return &csvEntry{ts: ts, fields: fields, headers: headers}, nil
}
