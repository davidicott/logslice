package parser

import "fmt"

// Format names supported by the registry.
const (
	FormatNginx    = "nginx"
	FormatApache   = "apache"
	FormatCombined = "combined"
	FormatSyslog   = "syslog"
	FormatJSON     = "json"
	FormatCSV      = "csv"
)

// factory is a function that constructs a Parser, optionally using options.
type factory func(opts ...string) Parser

var registry = map[string]factory{
	FormatNginx:    func(_ ...string) Parser { return NewNginxParser() },
	FormatApache:   func(_ ...string) Parser { return NewApacheParser() },
	FormatCombined: func(_ ...string) Parser { return NewCombinedParser() },
	FormatSyslog:   func(_ ...string) Parser { return NewSyslogParser() },
	FormatJSON: func(opts ...string) Parser {
		field := "time"
		if len(opts) > 0 && opts[0] != "" {
			field = opts[0]
		}
		return NewJSONParser(field)
	},
	FormatCSV: func(opts ...string) Parser {
		col := 0
		if len(opts) > 0 && opts[0] != "" {
			fmt.Sscanf(opts[0], "%d", &col)
		}
		return NewCSVParser(col)
	},
}

// NewFromFormat returns a Parser for the named format.
// For json, opts[0] may specify the timestamp field name.
// For csv, opts[0] may specify the timestamp column index as a string.
func NewFromFormat(format string, opts ...string) (Parser, error) {
	f, ok := registry[format]
	if !ok {
		return nil, fmt.Errorf("parser: unknown format %q", format)
	}
	return f(opts...), nil
}

// Formats returns the list of registered format names.
func Formats() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}
