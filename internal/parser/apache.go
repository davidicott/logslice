package parser

import (
	"fmt"
	"regexp"
	"time"
)

// ApacheCombinedLayout is the time layout used in Apache Combined Log Format.
const ApacheCombinedLayout = "02/Jan/2006:15:04:05 -0700"

// apacheCombinedRe matches Apache Combined Log Format lines.
// Example: 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326
var apacheCombinedRe = regexp.MustCompile(
	`^(\S+) \S+ (\S+) \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d{3}) (\d+|-)`,
)

// ApacheEntry represents a single parsed Apache access log entry.
type ApacheEntry struct {
	IP        string
	User      string
	Time      time.Time
	Method    string
	Path      string
	Protocol  string
	Status    int
	Bytes     int
	Raw       string
}

// Timestamp implements the LogEntry interface.
func (e *ApacheEntry) Timestamp() time.Time { return e.Time }

// ApacheParser parses Apache Combined Log Format lines.
type ApacheParser struct{}

// NewApacheParser returns a new ApacheParser.
func NewApacheParser() *ApacheParser {
	return &ApacheParser{}
}

// Parse parses a single Apache Combined Log Format line.
// Returns nil, nil for empty lines.
func (p *ApacheParser) Parse(line string) (LogEntry, error) {
	if line == "" {
		return nil, nil
	}

	matches := apacheCombinedRe.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("apache: line did not match combined log format: %q", line)
	}

	t, err := time.Parse(ApacheCombinedLayout, matches[3])
	if err != nil {
		return nil, fmt.Errorf("apache: failed to parse time %q: %w", matches[3], err)
	}

	var status, bytes int
	fmt.Sscanf(matches[7], "%d", &status)
	if matches[8] != "-" {
		fmt.Sscanf(matches[8], "%d", &bytes)
	}

	return &ApacheEntry{
		IP:       matches[1],
		User:     matches[2],
		Time:     t,
		Method:   matches[4],
		Path:     matches[5],
		Protocol: matches[6],
		Status:   status,
		Bytes:    bytes,
		Raw:      line,
	}, nil
}
