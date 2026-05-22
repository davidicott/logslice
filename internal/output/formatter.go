package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/example/logslice/internal/parser"
)

// Format represents the output serialization format.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// jsonRecord is the JSON representation of a log entry.
type jsonRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Raw       string    `json:"raw"`
}

// Formatter writes log entries to an io.Writer in the specified format.
type Formatter struct {
	w      io.Writer
	format Format
	enc    *json.Encoder
}

// NewFormatter creates a Formatter that writes to w using the given format.
func NewFormatter(w io.Writer, format Format) *Formatter {
	f := &Formatter{w: w, format: format}
	if format == FormatJSON {
		f.enc = json.NewEncoder(w)
	}
	return f
}

// Write serializes a single entry to the underlying writer.
func (f *Formatter) Write(entry parser.Entry) error {
	switch f.format {
	case FormatJSON:
		return f.enc.Encode(jsonRecord{
			Timestamp: entry.Timestamp(),
			Raw:       entry.Raw(),
		})
	default:
		_, err := fmt.Fprintln(f.w, entry.Raw())
		return err
	}
}

// WriteAll serializes all entries.
func (f *Formatter) WriteAll(entries []parser.Entry) error {
	for _, e := range entries {
		if err := f.Write(e); err != nil {
			return err
		}
	}
	return nil
}
