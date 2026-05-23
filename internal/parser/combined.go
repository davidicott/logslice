package parser

import (
	"fmt"
	"regexp"
	"time"
)

// CombinedEntry represents a parsed Apache/Nginx combined log format entry.
type CombinedEntry struct {
	IP        string
	Timestamp time.Time
	Method    string
	Path      string
	Status    int
	Bytes     int
	Referer   string
	UserAgent string
	Raw       string
}

func (e *CombinedEntry) Time() time.Time { return e.Timestamp }
func (e *CombinedEntry) String() string  { return e.Raw }

// combinedParser parses the "combined" log format used by both Apache and Nginx.
type combinedParser struct {
	re *regexp.Regexp
}

// NewCombinedParser returns a Parser for the combined log format:
// %h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i"
func NewCombinedParser() Parser {
	pattern := `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) \S+" (\d{3}) (\d+|-) "([^"]*)" "([^"]*)"`
	return &combinedParser{
		re: regexp.MustCompile(pattern),
	}
}

const combinedTimeLayout = "02/Jan/2006:15:04:05 -0700"

func (p *combinedParser) Parse(line string) (Entry, error) {
	if line == "" {
		return nil, ErrEmptyLine
	}
	m := p.re.FindStringSubmatch(line)
	if m == nil {
		return nil, fmt.Errorf("combined: no match: %q", line)
	}
	t, err := time.Parse(combinedTimeLayout, m[2])
	if err != nil {
		return nil, fmt.Errorf("combined: bad timestamp %q: %w", m[2], err)
	}
	status := 0
	fmt.Sscanf(m[5], "%d", &status)
	bytes := 0
	if m[6] != "-" {
		fmt.Sscanf(m[6], "%d", &bytes)
	}
	return &CombinedEntry{
		IP:        m[1],
		Timestamp: t,
		Method:    m[3],
		Path:      m[4],
		Status:    status,
		Bytes:     bytes,
		Referer:   m[7],
		UserAgent: m[8],
		Raw:       line,
	}, nil
}
