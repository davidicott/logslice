// Package parser provides log entry parsers for common web server formats.
package parser

import "time"

// LogEntry is the common interface implemented by all parsed log entries.
type LogEntry interface {
	// Timestamp returns the time at which the log event occurred.
	Timestamp() time.Time
}

// Parser is the common interface for log line parsers.
type Parser interface {
	// Parse parses a single log line and returns a LogEntry.
	// Returns (nil, nil) for empty or skipped lines.
	// Returns (nil, error) if the line cannot be parsed.
	Parse(line string) (LogEntry, error)
}
