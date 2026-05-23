package parser

import (
	"testing"
	"time"
)

const sampleCombined = `127.0.0.1 - frank [10/Oct/2023:13:55:36 -0700] "GET /index.html HTTP/1.1" 200 2326 "http://example.com/" "Mozilla/5.0"`

func TestCombinedParser_Parse_Valid(t *testing.T) {
	p := NewCombinedParser()
	entry, err := p.Parse(sampleCombined)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := entry.(*CombinedEntry)
	if !ok {
		t.Fatal("expected *CombinedEntry")
	}
	if e.IP != "127.0.0.1" {
		t.Errorf("IP: got %q, want %q", e.IP, "127.0.0.1")
	}
	if e.Method != "GET" {
		t.Errorf("Method: got %q, want %q", e.Method, "GET")
	}
	if e.Path != "/index.html" {
		t.Errorf("Path: got %q, want %q", e.Path, "/index.html")
	}
	if e.Status != 200 {
		t.Errorf("Status: got %d, want 200", e.Status)
	}
	if e.Bytes != 2326 {
		t.Errorf("Bytes: got %d, want 2326", e.Bytes)
	}
	if e.Referer != "http://example.com/" {
		t.Errorf("Referer: got %q", e.Referer)
	}
	if e.UserAgent != "Mozilla/5.0" {
		t.Errorf("UserAgent: got %q", e.UserAgent)
	}
}

func TestCombinedParser_Parse_EmptyLine(t *testing.T) {
	p := NewCombinedParser()
	_, err := p.Parse("")
	if err != ErrEmptyLine {
		t.Errorf("expected ErrEmptyLine, got %v", err)
	}
}

func TestCombinedParser_Parse_InvalidLine(t *testing.T) {
	p := NewCombinedParser()
	_, err := p.Parse("not a log line at all")
	if err == nil {
		t.Error("expected error for invalid line")
	}
}

func TestCombinedParser_Parse_DashBytes(t *testing.T) {
	line := `10.0.0.1 - - [10/Oct/2023:14:00:00 +0000] "POST /api HTTP/1.1" 204 - "-" "curl/7.68"`
	p := NewCombinedParser()
	entry, err := p.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := entry.(*CombinedEntry)
	if e.Bytes != 0 {
		t.Errorf("Bytes: got %d, want 0 for dash", e.Bytes)
	}
	if e.Status != 204 {
		t.Errorf("Status: got %d, want 204", e.Status)
	}
}

func TestCombinedEntry_Timestamp(t *testing.T) {
	p := NewCombinedParser()
	entry, err := p.Parse(sampleCombined)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2023, time.October, 10, 13, 55, 36, 0, entry.Time().Location())
	if !entry.Time().Equal(want) {
		t.Errorf("Time: got %v, want %v", entry.Time(), want)
	}
}
