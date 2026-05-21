package parser

import (
	"testing"
	"time"
)

const sampleNginxLine = `127.0.0.1 - frank [10/Oct/2023:13:55:36 -0700] "GET /index.html HTTP/1.1" 200 2326 "http://example.com/" "Mozilla/5.0"`

func TestNginxParser_Parse_Valid(t *testing.T) {
	p := NewNginxParser()
	entry, err := p.Parse(sampleNginxLine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entry.IP != "127.0.0.1" {
		t.Errorf("IP: got %q, want %q", entry.IP, "127.0.0.1")
	}
	if entry.Method != "GET" {
		t.Errorf("Method: got %q, want %q", entry.Method, "GET")
	}
	if entry.Path != "/index.html" {
		t.Errorf("Path: got %q, want %q", entry.Path, "/index.html")
	}
	if entry.Protocol != "HTTP/1.1" {
		t.Errorf("Protocol: got %q, want %q", entry.Protocol, "HTTP/1.1")
	}
	if entry.Status != 200 {
		t.Errorf("Status: got %d, want %d", entry.Status, 200)
	}
	if entry.BodyBytes != 2326 {
		t.Errorf("BodyBytes: got %d, want %d", entry.BodyBytes, 2326)
	}
	if entry.Referer != "http://example.com/" {
		t.Errorf("Referer: got %q, want %q", entry.Referer, "http://example.com/")
	}

	wantTime := time.Date(2023, time.October, 10, 13, 55, 36, 0, entry.Time.Location())
	if !entry.Time.Equal(wantTime) {
		t.Errorf("Time: got %v, want %v", entry.Time, wantTime)
	}
	if entry.Raw != sampleNginxLine {
		t.Errorf("Raw: got %q, want original line", entry.Raw)
	}
}

func TestNginxParser_Parse_InvalidLine(t *testing.T) {
	p := NewNginxParser()
	_, err := p.Parse("this is not a valid nginx log line")
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestNginxParser_Parse_EmptyLine(t *testing.T) {
	p := NewNginxParser()
	_, err := p.Parse("")
	if err == nil {
		t.Fatal("expected error for empty line, got nil")
	}
}

func TestNginxEntry_Timestamp(t *testing.T) {
	p := NewNginxParser()
	entry, err := p.Parse(sampleNginxLine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp() != entry.Time {
		t.Errorf("Timestamp() = %v, want %v", entry.Timestamp(), entry.Time)
	}
}
