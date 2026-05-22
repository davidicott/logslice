package filter

import (
	"github.com/example/logslice/internal/parser"
	"github.com/example/logslice/internal/timerange"
)

// Filter applies a time range to a stream of log entries.
type Filter struct {
	tr *timerange.TimeRange
}

// New creates a new Filter for the given TimeRange.
func New(tr *timerange.TimeRange) *Filter {
	return &Filter{tr: tr}
}

// Match returns true if the entry's timestamp falls within the filter's time range.
func (f *Filter) Match(entry parser.Entry) bool {
	t := entry.Timestamp()
	return f.tr.Contains(t)
}

// Apply filters a slice of entries, returning only those within the time range.
func (f *Filter) Apply(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if f.Match(e) {
			result = append(result, e)
		}
	}
	return result
}
