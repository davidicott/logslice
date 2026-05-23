package parser

import (
	"fmt"
	"strings"
	"time"
)

// logfmt key=value pair log format
// Example: ts=2024-01-15T10:30:00Z level=info msg="user logged in" user=alice

type LogfmtEntry struct {
	ts     time.Time
	fields map[string]string
	raw    string
}

func (e *LogfmtEntry) Timestamp() time.Time { return e.ts }
func (e *LogfmtEntry) String() string        { return e.raw }

type LogfmtParser struct {
	timeKey string
	timeFmt string
}

func NewLogfmtParser(timeKey, timeFmt string) *LogfmtParser {
	if timeKey == "" {
		timeKey = "ts"
	}
	if timeFmt == "" {
		timeFmt = time.RFC3339
	}
	return &LogfmtParser{timeKey: timeKey, timeFmt: timeFmt}
}

func (p *LogfmtParser) Parse(line string) (Entry, error) {
	if strings.TrimSpace(line) == "" {
		return nil, ErrEmptyLine
	}

	fields, err := parseLogfmtFields(line)
	if err != nil {
		return nil, fmt.Errorf("logfmt parse error: %w", err)
	}

	tsVal, ok := fields[p.timeKey]
	if !ok {
		return nil, fmt.Errorf("logfmt: time key %q not found", p.timeKey)
	}

	t, err := time.Parse(p.timeFmt, tsVal)
	if err != nil {
		return nil, fmt.Errorf("logfmt: invalid timestamp %q: %w", tsVal, err)
	}

	return &LogfmtEntry{ts: t, fields: fields, raw: line}, nil
}

func parseLogfmtFields(line string) (map[string]string, error) {
	fields := make(map[string]string)
	remaining := strings.TrimSpace(line)

	for remaining != "" {
		eqIdx := strings.IndexByte(remaining, '=')
		if eqIdx <= 0 {
			return nil, fmt.Errorf("malformed key=value near %q", remaining)
		}
		key := remaining[:eqIdx]
		remaining = remaining[eqIdx+1:]

		var val string
		if strings.HasPrefix(remaining, `"`) {
			closing := strings.Index(remaining[1:], `"`)
			if closing < 0 {
				return nil, fmt.Errorf("unterminated quoted value for key %q", key)
			}
			val = remaining[1 : closing+1]
			remaining = strings.TrimSpace(remaining[closing+2:])
		} else {
			spIdx := strings.IndexByte(remaining, ' ')
			if spIdx < 0 {
				val = remaining
				remaining = ""
			} else {
				val = remaining[:spIdx]
				remaining = strings.TrimSpace(remaining[spIdx+1:])
			}
		}
		fields[key] = val
	}
	return fields, nil
}
