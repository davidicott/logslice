package parser

import (
	"testing"
	"time"
)

const validApacheLine = `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
const noBodyApacheLine = `192.168.1.1 - - [01/Jan/2024:00:00:01 +0000] "HEAD / HTTP/1.1" 204 -`

func TestApacheParser_Parse_Valid(t *testing.T) {
	p := NewApacheParser()
	entry, err := p.Parse(validApacheLine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}

	a, ok := entry.(*ApacheEntry)
	if !ok {
		t.Fatalf("expected *ApacheEntry, got %T", entry)
	}

	if a.IP != "127.0.0.1" {
		t.Errorf("IP: got %q, want %q", a.IP, "127.0.0.1")
	}
	if a.User != "frank" {
		t.Errorf("User: got %q, want %q", a.User, "frank")
	}
	if a.Method != "GET" {
		t.Errorf("Method: got %q, want %q", a.Method, "GET")
	}
	if a.Path != "/apache_pb.gif" {
		t.Errorf("Path: got %q, want %q", a.Path, "/apache_pb.gif")
	}
	if a.Status != 200 {
		t.Errorf("Status: got %d, want 200", a.Status)
	}
	if a.Bytes != 2326 {
		t.Errorf("Bytes: got %d, want 2326", a.Bytes)
	}
}

func TestApacheParser_Parse_NoDash(t *testing.T) {
	p := NewApacheParser()
	entry, err := p.Parse(noBodyApacheLine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := entry.(*ApacheEntry)
	if a.Bytes != 0 {
		t.Errorf("Bytes: got %d, want 0 for dash", a.Bytes)
	}
	if a.Status != 204 {
		t.Errorf("Status: got %d, want 204", a.Status)
	}
}

func TestApacheParser_Parse_EmptyLine(t *testing.T) {
	p := NewApacheParser()
	entry, err := p.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil entry for empty line")
	}
}

func TestApacheParser_Parse_InvalidLine(t *testing.T) {
	p := NewApacheParser()
	_, err := p.Parse("not a valid log line")
	if err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestApacheEntry_Timestamp(t *testing.T) {
	p := NewApacheParser()
	entry, err := p.Parse(validApacheLine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, _ := time.Parse(ApacheCombinedLayout, "10/Oct/2000:13:55:36 -0700")
	if !entry.Timestamp().Equal(want) {
		t.Errorf("Timestamp: got %v, want %v", entry.Timestamp(), want)
	}
}
