package parser

import (
	"testing"
	"time"
)

func TestJSONParser_Parse_Valid(t *testing.T) {
	p := NewJSONParser("", "")
	line := `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"started"}`

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entry.GetTimestamp().Equal(expected) {
		t.Errorf("got timestamp %v, want %v", entry.GetTimestamp(), expected)
	}
}

func TestJSONParser_Parse_CustomField(t *testing.T) {
	p := NewJSONParser("@timestamp", time.RFC3339)
	line := `{"@timestamp":"2024-03-01T08:30:00Z","service":"api"}`

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2024, 3, 1, 8, 30, 0, 0, time.UTC)
	if !entry.GetTimestamp().Equal(expected) {
		t.Errorf("got timestamp %v, want %v", entry.GetTimestamp(), expected)
	}
}

func TestJSONParser_Parse_EmptyLine(t *testing.T) {
	p := NewJSONParser("", "")
	_, err := p.Parse("")
	if err != ErrEmptyLine {
		t.Errorf("expected ErrEmptyLine, got %v", err)
	}
}

func TestJSONParser_Parse_InvalidJSON(t *testing.T) {
	p := NewJSONParser("", "")
	_, err := p.Parse("not json at all")
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestJSONParser_Parse_MissingTimeField(t *testing.T) {
	p := NewJSONParser("", "")
	_, err := p.Parse(`{"level":"warn","msg":"no timestamp here"}`)
	if err == nil {
		t.Error("expected error for missing time field, got nil")
	}
}

func TestJSONParser_Parse_BadTimestamp(t *testing.T) {
	p := NewJSONParser("", "")
	_, err := p.Parse(`{"time":"not-a-date","msg":"bad"}`)
	if err == nil {
		t.Error("expected error for unparseable timestamp, got nil")
	}
}

func TestJSONEntry_String(t *testing.T) {
	p := NewJSONParser("", "")
	line := `{"time":"2024-01-15T10:00:00Z","msg":"hello"}`

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := entry.String()
	if s == "" {
		t.Error("expected non-empty string from JSONEntry.String()")
	}
}
