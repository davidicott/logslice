package parser

import (
	"fmt"
	"regexp"
	"time"
)

// NginxEntry represents a parsed nginx combined log format entry.
type NginxEntry struct {
	IP        string
	Time      time.Time
	Method    string
	Path      string
	Protocol  string
	Status    int
	BodyBytes int
	Referer   string
	UserAgent string
	Raw       string
}

// nginxCombinedRegex matches the nginx combined log format:
// $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
var nginxCombinedRegex = regexp.MustCompile(
	`^(\S+) - \S+ \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d+) (\d+) "([^"]*)" "([^"]*)"`,
)

const nginxTimeLayout = "02/Jan/2006:15:04:05 -0700"

// NginxParser parses nginx combined log format lines.
type NginxParser struct{}

// NewNginxParser creates a new NginxParser.
func NewNginxParser() *NginxParser {
	return &NginxParser{}
}

// Parse parses a single nginx log line and returns an NginxEntry.
func (p *NginxParser) Parse(line string) (*NginxEntry, error) {
	matches := nginxCombinedRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("line does not match nginx combined log format")
	}

	t, err := time.Parse(nginxTimeLayout, matches[2])
	if err != nil {
		return nil, fmt.Errorf("parse time %q: %w", matches[2], err)
	}

	var status, bodyBytes int
	fmt.Sscanf(matches[6], "%d", &status)
	fmt.Sscanf(matches[7], "%d", &bodyBytes)

	return &NginxEntry{
		IP:        matches[1],
		Time:      t,
		Method:    matches[3],
		Path:      matches[4],
		Protocol:  matches[5],
		Status:    status,
		BodyBytes: bodyBytes,
		Referer:   matches[8],
		UserAgent: matches[9],
		Raw:       line,
	}, nil
}

// Timestamp returns the log entry's timestamp, satisfying a common Entry interface.
func (e *NginxEntry) Timestamp() time.Time {
	return e.Time
}
