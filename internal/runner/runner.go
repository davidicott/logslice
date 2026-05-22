// Package runner wires together parsing, filtering, and formatting
// into a single reusable pipeline.
package runner

import (
	"bufio"
	"fmt"
	"io"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/parser"
)

// Run reads lines from r, parses each with p, filters with f, and writes
// matching entries to fmt. It returns the number of matched entries and any
// non-EOF error encountered.
func Run(r io.Reader, p parser.Parser, f *filter.Filter, fmt *output.Formatter) (int, error) {
	scanner := bufio.NewScanner(r)
	var entries []parser.Entry

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, err := p.Parse(line)
		if err != nil {
			// Skip unparseable lines; non-fatal.
			continue
		}

		if f.Match(entry) {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	if err := fmt.WriteAll(entries); err != nil {
		return 0, fmt.Errorf("write error: %w", err)
	}

	return len(entries), nil
}

// Errorf is a small helper so runner does not depend on the fmt package
// at the call site — kept internal to avoid confusion with output.Formatter.
func errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
