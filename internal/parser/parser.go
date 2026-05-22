package parser

import "time"

// Entry represents a single parsed log entry.
type Entry interface {
	// Timestamp returns the time associated with this log entry.
	Timestamp() time.Time
	// String returns the original raw log line.
	String() string
}

// Parser parses a single log line into an Entry.
type Parser interface {
	Parse(line string) (Entry, error)
}
