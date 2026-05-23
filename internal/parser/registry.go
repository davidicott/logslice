package parser

import "fmt"

// NewFromFormat returns a Parser for the given named format.
// Supported formats: nginx, apache, combined, syslog, json, csv, logfmt, w3c.
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
		return NewJSONParser(""), nil
	case "csv":
		return NewCSVParser("", 0), nil
	case "logfmt":
		return NewLogfmtParser(""), nil
	case "w3c":
		return NewW3CParser(), nil
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}
