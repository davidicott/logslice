package parser

import (
	"fmt"
	"strings"
	"time"
)

// W3C Extended Log Format parser
// Handles logs with #Fields directive header
// https://www.w3.org/TR/WD-logfile

const w3cTimeLayout = "2006-01-02 15:04:05"

type w3cEntry struct {
	raw       string
	timestamp time.Time
}

func (e *w3cEntry) Timestamp() time.Time { return e.timestamp }
func (e *w3cEntry) String() string        { return e.raw }

type w3cParser struct {
	dateIdx int
	timeIdx int
}

// NewW3CParser returns a Parser for W3C Extended Log Format.
// It reads the #Fields directive to locate date and time columns.
func NewW3CParser() Parser {
	return &w3cParser{dateIdx: -1, timeIdx: -1}
}

func (p *w3cParser) Parse(line string) (Entry, error) {
	if line == "" {
		return nil, ErrEmptyLine
	}

	// Handle directive lines
	if strings.HasPrefix(line, "#") {
		if strings.HasPrefix(line, "#Fields:") {
			p.parseFields(line)
		}
		return nil, ErrEmptyLine
	}

	if p.dateIdx == -1 || p.timeIdx == -1 {
		return nil, fmt.Errorf("w3c: no #Fields directive seen yet")
	}

	fields := strings.Fields(line)
	if p.dateIdx >= len(fields) || p.timeIdx >= len(fields) {
		return nil, fmt.Errorf("w3c: not enough fields in line")
	}

	datePart := fields[p.dateIdx]
	timePart := fields[p.timeIdx]

	ts, err := time.Parse(w3cTimeLayout, datePart+" "+timePart)
	if err != nil {
		return nil, fmt.Errorf("w3c: parse time %q %q: %w", datePart, timePart, err)
	}

	return &w3cEntry{raw: line, timestamp: ts}, nil
}

func (p *w3cParser) parseFields(line string) {
	// "#Fields: date time cs-method cs-uri ..."
	parts := strings.Fields(strings.TrimPrefix(line, "#Fields:"))
	for i, f := range parts {
		switch f {
		case "date":
			p.dateIdx = i
		case "time":
			p.timeIdx = i
		}
	}
}
