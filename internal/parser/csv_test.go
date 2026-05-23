package parser

import (
	"testing"
	"time"
)

func TestCSVParser_Parse_HeaderThenData(t *testing.T) {
	p := NewCSVParser(0, time.RFC3339)

	// First line should be treated as header — returns nil, nil.
	entry, err := p.Parse("timestamp,level,message")
	if err != nil {
		t.Fatalf("unexpected error on header: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil entry for header row")
	}

	// Data line.
	entry, err = p.Parse("2024-01-15T10:00:00Z,INFO,hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}

	want := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entry.Timestamp().Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp(), want)
	}
}

func TestCSVParser_Parse_EmptyLine(t *testing.T) {
	p := NewCSVParser(0, time.RFC3339)
	entry, err := p.Parse("   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil entry for empty line")
	}
}

func TestCSVParser_Parse_InvalidTimestamp(t *testing.T) {
	p := NewCSVParser(0, time.RFC3339)
	p.headerParsed = true // skip header phase
	_, err := p.Parse("not-a-time,INFO,msg")
	if err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}

func TestCSVParser_Parse_TimeColumnOutOfRange(t *testing.T) {
	p := NewCSVParser(5, time.RFC3339)
	p.headerParsed = true
	_, err := p.Parse("2024-01-15T10:00:00Z,INFO")
	if err == nil {
		t.Fatal("expected error when time column is out of range")
	}
}

func TestCSVEntry_String(t *testing.T) {
	p := NewCSVParser(0, time.RFC3339)
	// Feed header.
	p.Parse("timestamp,level,message") //nolint:errcheck

	entry, err := p.Parse("2024-03-01T08:30:00Z,WARN,disk full")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := entry.String()
	want := "timestamp=2024-03-01T08:30:00Z level=WARN message=disk full"
	if got != want {
		t.Errorf("String(): got %q, want %q", got, want)
	}
}

func TestCSVParser_NoHeaderFallback(t *testing.T) {
	p := NewCSVParser(0, time.RFC3339)
	p.headerParsed = true // skip header, no headers set

	entry, err := p.Parse("2024-06-01T00:00:00Z,DEBUG,test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := entry.String()
	if got == "" {
		t.Error("expected non-empty string output")
	}
}
