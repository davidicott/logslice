package parser

import (
	"testing"
	"time"
)

func TestW3CParser_Parse_Valid(t *testing.T) {
	p := NewW3CParser()

	// Feed the fields directive first
	_, err := p.Parse("#Fields: date time cs-method cs-uri sc-status")
	if err != ErrEmptyLine {
		t.Fatalf("expected ErrEmptyLine for directive, got %v", err)
	}

	line := "2024-03-15 12:34:56 GET /index.html 200"
	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(2024, 3, 15, 12, 34, 56, 0, time.UTC)
	if !entry.Timestamp().Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp(), want)
	}
	if entry.String() != line {
		t.Errorf("String(): got %q, want %q", entry.String(), line)
	}
}

func TestW3CParser_Parse_EmptyLine(t *testing.T) {
	p := NewW3CParser()
	_, err := p.Parse("")
	if err != ErrEmptyLine {
		t.Errorf("expected ErrEmptyLine, got %v", err)
	}
}

func TestW3CParser_Parse_CommentLine(t *testing.T) {
	p := NewW3CParser()
	_, err := p.Parse("#Version: 1.0")
	if err != ErrEmptyLine {
		t.Errorf("expected ErrEmptyLine for comment, got %v", err)
	}
}

func TestW3CParser_Parse_NoFieldsDirective(t *testing.T) {
	p := NewW3CParser()
	_, err := p.Parse("2024-03-15 12:34:56 GET /index.html 200")
	if err == nil {
		t.Error("expected error when no #Fields directive seen")
	}
}

func TestW3CParser_Parse_InvalidTimestamp(t *testing.T) {
	p := NewW3CParser()
	p.Parse("#Fields: date time cs-method") //nolint

	_, err := p.Parse("not-a-date not-a-time GET")
	if err == nil {
		t.Error("expected error for invalid timestamp")
	}
}

func TestW3CParser_Parse_FieldsOrderVariant(t *testing.T) {
	p := NewW3CParser()
	p.Parse("#Fields: cs-method time date sc-status") //nolint

	line := "GET 08:00:00 2024-06-01 200"
	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	if !entry.Timestamp().Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp(), want)
	}
}

func TestW3CParser_Parse_NotEnoughFields(t *testing.T) {
	p := NewW3CParser()
	p.Parse("#Fields: date time cs-method") //nolint

	_, err := p.Parse("2024-03-15") // only one field, need at least 2
	if err == nil {
		t.Error("expected error for not enough fields")
	}
}
