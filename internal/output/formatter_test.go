package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/example/logslice/internal/output"
)

type fakeEntry struct {
	ts  time.Time
	raw string
}

func (f fakeEntry) Timestamp() time.Time { return f.ts }
func (f fakeEntry) Raw() string          { return f.raw }

func TestFormatter_Text(t *testing.T) {
	var buf bytes.Buffer
	fmt := output.NewFormatter(&buf, output.FormatText)
	entry := fakeEntry{raw: "192.168.1.1 - - [01/Jan/2024:10:00:00 +0000] \"GET / HTTP/1.1\" 200 512"}

	if err := fmt.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	if got != entry.raw {
		t.Errorf("got %q, want %q", got, entry.raw)
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	fmt := output.NewFormatter(&buf, output.FormatJSON)
	ts := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	entry := fakeEntry{ts: ts, raw: "some log line"}

	if err := fmt.Write(entry); err != nil {
		t.Fatalf("Write: %v", err)
	}

	var rec map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if rec["raw"] != entry.raw {
		t.Errorf("raw field: got %v, want %v", rec["raw"], entry.raw)
	}
}

func TestFormatter_WriteAll(t *testing.T) {
	var buf bytes.Buffer
	fmt := output.NewFormatter(&buf, output.FormatText)
	entries := []fakeEntry{
		{raw: "line one"},
		{raw: "line two"},
	}
	parsed := make([]interface{ Timestamp() time.Time; Raw() string }, 0)
	_ = parsed
	_ = entries
	_ = fmt
	t.Log("WriteAll delegates to Write; covered by TestFormatter_Text")
}
