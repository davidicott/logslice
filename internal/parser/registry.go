package parser

import (
	"fmt"
	"time"
)

// Formats lists all supported log format names.
var Formats = []string{"nginx", "apache", "combined", "syslog", "json", "csv", "logfmt"}

// NewFromFormat returns a Parser for the given format name using sensible defaults.
// Returns an error if the format is unknown.
func NewFromFormat(format string) (Parser, error) {
	switch format {
	case "nginx":
		return NewNginxParser(), nil
	case "apache":
		return NewApacheParser(), nil
	case "combined":
		return NewCombinedParser(), nil
	case "syslog":
		return NewSyslogParser(), nil
	case "json":
		return NewJSONParser("", ""), nil
	case "csv":
		return NewCSVParser(0, time.RFC3339), nil
	case "logfmt":
		return NewLogfmtParser("", ""), nil
	default:
		return nil, fmt.Errorf("unknown format %q: supported formats are %v", format, Formats)
	}
}
