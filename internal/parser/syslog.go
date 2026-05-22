package parser

import (
	"fmt"
	"regexp"
	"time"
)

// syslogLayout is the standard syslog timestamp format (RFC3164).
const syslogLayout = "Jan  2 15:04:05"
const syslogLayoutPadded = "Jan 02 15:04:05"

// syslogRe matches lines like: Oct  5 12:34:56 hostname process[pid]: message
var syslogRe = regexp.MustCompile(
	`^(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+(\S+)\s+(\S+):\s+(.*)$`,
)

// SyslogEntry represents a parsed syslog log line.
type SyslogEntry struct {
	Time     time.Time
	Host     string
	Process  string
	Message  string
	Raw      string
}

// Timestamp implements Entry.
func (e *SyslogEntry) Timestamp() time.Time { return e.Time }

// String implements Entry.
func (e *SyslogEntry) String() string { return e.Raw }

// syslogParser parses syslog-formatted log lines.
type syslogParser struct {
	year int
}

// NewSyslogParser returns a Parser for syslog (RFC3164) format.
// year is used to fill in the missing year in syslog timestamps.
func NewSyslogParser(year int) Parser {
	return &syslogParser{year: year}
}

// Parse parses a single syslog line.
func (p *syslogParser) Parse(line string) (Entry, error) {
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}
	matches := syslogRe.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("syslog: line does not match expected format: %q", line)
	}

	rawTime := matches[1]
	t, err := time.Parse(syslogLayout, rawTime)
	if err != nil {
		t, err = time.Parse(syslogLayoutPadded, rawTime)
		if err != nil {
			return nil, fmt.Errorf("syslog: cannot parse timestamp %q: %w", rawTime, err)
		}
	}
	// Syslog has no year; inject the provided year in UTC.
	t = time.Date(p.year, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)

	return &SyslogEntry{
		Time:    t,
		Host:    matches[2],
		Process: matches[3],
		Message: matches[4],
		Raw:     line,
	}, nil
}
