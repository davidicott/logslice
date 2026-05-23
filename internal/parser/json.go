package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONEntry represents a parsed JSON log entry.
type JSONEntry struct {
	Raw       map[string]interface{}
	Timestamp time.Time
}

// Timestamp returns the entry's timestamp.
func (e *JSONEntry) GetTimestamp() time.Time {
	return e.Timestamp
}

// String returns a JSON representation of the entry.
func (e *JSONEntry) String() string {
	b, err := json.Marshal(e.Raw)
	if err != nil {
		return ""
	}
	return string(b)
}

// JSONParser parses JSON log lines where each line is a JSON object
// containing a timestamp field.
type JSONParser struct {
	timeField  string
	timeFormat string
}

// NewJSONParser creates a new JSONParser. timeField is the key used for the
// timestamp, and timeFormat is the Go time layout string.
func NewJSONParser(timeField, timeFormat string) *JSONParser {
	if timeField == "" {
		timeField = "time"
	}
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	return &JSONParser{timeField: timeField, timeFormat: timeFormat}
}

// Parse parses a single JSON log line and returns an Entry.
func (p *JSONParser) Parse(line string) (Entry, error) {
	if line == "" {
		return nil, ErrEmptyLine
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}

	v, ok := raw[p.timeField]
	if !ok {
		return nil, fmt.Errorf("missing time field %q", p.timeField)
	}

	ts, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("time field %q is not a string", p.timeField)
	}

	t, err := time.Parse(p.timeFormat, ts)
	if err != nil {
		return nil, fmt.Errorf("cannot parse timestamp %q: %w", ts, err)
	}

	return &JSONEntry{Raw: raw, Timestamp: t}, nil
}
