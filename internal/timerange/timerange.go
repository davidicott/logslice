package timerange

import (
	"fmt"
	"time"
)

const timeLayout = time.RFC3339

// TimeRange represents an inclusive start/end time window.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// New creates a TimeRange from RFC3339 string values.
func New(start, end string) (*TimeRange, error) {
	s, err := time.Parse(timeLayout, start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time %q: %w", start, err)
	}
	e, err := time.Parse(timeLayout, end)
	if err != nil {
		return nil, fmt.Errorf("invalid end time %q: %w", end, err)
	}
	if !s.Before(e) {
		return nil, fmt.Errorf("start time must be before end time")
	}
	return &TimeRange{Start: s, End: e}, nil
}

// Parse creates a TimeRange from two already-parsed time.Time values.
func Parse(start, end time.Time) (*TimeRange, error) {
	if !start.Before(end) {
		return nil, fmt.Errorf("start time must be before end time")
	}
	return &TimeRange{Start: start, End: end}, nil
}

// Contains reports whether t falls within the inclusive range [Start, End].
func (tr *TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && !t.After(tr.End)
}

// String returns a human-readable representation of the TimeRange.
func (tr *TimeRange) String() string {
	return fmt.Sprintf("[%s, %s]", tr.Start.Format(timeLayout), tr.End.Format(timeLayout))
}
