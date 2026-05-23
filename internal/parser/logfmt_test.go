package parser

import (
	"testing"
	"time"
)

func TestLogfmtParser_Parse_Valid(t *testing.T) {
	p := NewLogfmtParser("", "")
	line := `ts=2024-01-15T10:30:00Z level=info msg="user logged in" user=alice`

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !entry.Timestamp().Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp(), want)
	}
	if entry.String() != line {
		t.Errorf("String(): got %q, want %q", entry.String(), line)
	}
}

func TestLogfmtParser_Parse_CustomKey(t *testing.T) {
	p := NewLogfmtParser("time", time.RFC3339)
	line := `time=2024-03-01T08:00:00Z level=warn msg=slowquery`

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, 3, 1, 8, 0, 0, 0, time.UTC)
	if !entry.Timestamp().Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp(), want)
	}
}

func TestLogfmtParser_Parse_EmptyLine(t *testing.T) {
	p := NewLogfmtParser("", "")
	_, err := p.Parse("")
	if err != ErrEmptyLine {
		t.Errorf("expected ErrEmptyLine, got %v", err)
	}
}

func TestLogfmtParser_Parse_MissingTimeKey(t *testing.T) {
	p := NewLogfmtParser("ts", "")
	line := `level=info msg=hello`
	_, err := p.Parse(line)
	if err == nil {
		t.Error("expected error for missing time key, got nil")
	}
}

func TestLogfmtParser_Parse_InvalidTimestamp(t *testing.T) {
	p := NewLogfmtParser("", "")
	line := `ts=not-a-time level=info`
	_, err := p.Parse(line)
	if err == nil {
		t.Error("expected error for invalid timestamp, got nil")
	}
}

func TestLogfmtParser_Parse_MalformedLine(t *testing.T) {
	p := NewLogfmtParser("", "")
	line := `justakeywithnoequalssign`
	_, err := p.Parse(line)
	if err == nil {
		t.Error("expected error for malformed line, got nil")
	}
}

func TestLogfmtEntry_String(t *testing.T) {
	p := NewLogfmtParser("", "")
	line := `ts=2024-06-01T12:00:00Z msg=test`
	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.String() != line {
		t.Errorf("String(): got %q, want %q", entry.String(), line)
	}
}
