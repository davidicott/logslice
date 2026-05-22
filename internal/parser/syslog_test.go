package parser

import (
	"testing"
	"time"
)

const testYear = 2024

func TestSyslogParser_Parse_Valid(t *testing.T) {
	p := NewSyslogParser(testYear)
	line := "Oct  5 12:34:56 myhost myapp[1234]: something happened"

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	se, ok := entry.(*SyslogEntry)
	if !ok {
		t.Fatalf("expected *SyslogEntry, got %T", entry)
	}

	if se.Host != "myhost" {
		t.Errorf("host: got %q, want %q", se.Host, "myhost")
	}
	if se.Process != "myapp[1234]" {
		t.Errorf("process: got %q, want %q", se.Process, "myapp[1234]")
	}
	if se.Message != "something happened" {
		t.Errorf("message: got %q, want %q", se.Message, "something happened")
	}
	if se.Raw != line {
		t.Errorf("raw: got %q, want %q", se.Raw, line)
	}
}

func TestSyslogParser_Parse_EmptyLine(t *testing.T) {
	p := NewSyslogParser(testYear)
	_, err := p.Parse("")
	if err == nil {
		t.Fatal("expected error for empty line, got nil")
	}
}

func TestSyslogParser_Parse_InvalidLine(t *testing.T) {
	p := NewSyslogParser(testYear)
	_, err := p.Parse("not a syslog line at all")
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestSyslogEntry_Timestamp(t *testing.T) {
	p := NewSyslogParser(testYear)
	line := "Jan  2 15:04:05 host proc: msg"

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := time.Date(testYear, time.January, 2, 15, 4, 5, 0, time.UTC)
	if !entry.Timestamp().Equal(want) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp(), want)
	}
}

func TestSyslogEntry_String(t *testing.T) {
	p := NewSyslogParser(testYear)
	line := "Mar 15 08:00:00 server kernel: boot complete"

	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.String() != line {
		t.Errorf("String(): got %q, want %q", entry.String(), line)
	}
}
