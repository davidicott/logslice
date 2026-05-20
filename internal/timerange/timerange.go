package timerange

import (
	"fmt"
	"time"
)

// TimeRange represents an inclusive start/end window for log filtering.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// New creates a TimeRange from two time.Time values.
// Returns an error if start is after end.
func New(start, end time.Time) (TimeRange, error) {
	if start.After(end) {
		return TimeRange{}, fmt.Errorf("timerange: start %v is after end %v", start, end)
	}
	return TimeRange{Start: start, End: end}, nil
}

// Parse creates a TimeRange by parsing two RFC3339 timestamp strings.
func Parse(startStr, endStr string) (TimeRange, error) {
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return TimeRange{}, fmt.Errorf("timerange: invalid start time %q: %w", startStr, err)
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return TimeRange{}, fmt.Errorf("timerange: invalid end time %q: %w", endStr, err)
	}
	return New(start, end)
}

// Contains reports whether t falls within the inclusive range [Start, End].
func (tr TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}

// IsZero reports whether the TimeRange is the zero value.
func (tr TimeRange) IsZero() bool {
	return tr.Start.IsZero() && tr.End.IsZero()
}

// String returns a human-readable representation of the range.
func (tr TimeRange) String() string {
	return fmt.Sprintf("[%s, %s]", tr.Start.Format(time.RFC3339), tr.End.Format(time.RFC3339))
}
