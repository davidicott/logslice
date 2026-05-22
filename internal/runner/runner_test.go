package runner_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"logslice/internal/filter"
	"logslice/internal/output"
	"logslice/internal/parser"
	"logslice/internal/runner"
	"logslice/internal/timerange"
)

func mustRange(t *testing.T, start, end string) timerange.TimeRange {
	t.Helper()
	tr, err := timerange.New(start, end)
	if err != nil {
		t.Fatalf("mustRange: %v", err)
	}
	return tr
}

func TestRun_MatchesInRange(t *testing.T) {
	lines := strings.Join([]string{
		`127.0.0.1 - - [01/Jan/2024:10:00:00 +0000] "GET /a HTTP/1.1" 200 512`,
		`127.0.0.1 - - [01/Jan/2024:12:00:00 +0000] "GET /b HTTP/1.1" 404 128`,
		`127.0.0.1 - - [01/Jan/2024:14:00:00 +0000] "GET /c HTTP/1.1" 500 64`,
	}, "\n")

	tr := mustRange(t, "2024-01-01T11:00:00Z", "2024-01-01T13:00:00Z")
	f := filter.New(tr)
	var buf bytes.Buffer
	fmt := output.NewFormatter("text", &buf)

	n, err := runner.Run(strings.NewReader(lines), parser.NewApacheParser(), f, fmt)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 match, got %d", n)
	}
	if !strings.Contains(buf.String(), "/b") {
		t.Errorf("expected /b in output, got: %s", buf.String())
	}
}

func TestRun_NoMatches(t *testing.T) {
	lines := `127.0.0.1 - - [01/Jan/2024:08:00:00 +0000] "GET /x HTTP/1.1" 200 100`

	tr := mustRange(t, "2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")
	f := filter.New(tr)
	var buf bytes.Buffer
	fmt := output.NewFormatter("text", &buf)

	n, err := runner.Run(strings.NewReader(lines), parser.NewApacheParser(), f, fmt)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 matches, got %d", n)
	}
}

func TestRun_SkipsUnparseable(t *testing.T) {
	lines := strings.Join([]string{
		"this is not a valid log line",
		`127.0.0.1 - - [01/Jan/2024:11:30:00 +0000] "GET /ok HTTP/1.1" 200 200`,
	}, "\n")

	tr := mustRange(t, "2024-01-01T11:00:00Z", "2024-01-01T12:00:00Z")
	f := filter.New(tr)
	var buf bytes.Buffer
	fmt := output.NewFormatter("text", &buf)

	n, err := runner.Run(strings.NewReader(lines), parser.NewApacheParser(), f, fmt)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 match, got %d", n)
	}
}

// Ensure Run handles empty input gracefully.
func TestRun_EmptyInput(t *testing.T) {
	tr := mustRange(t, "2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z")
	f := filter.New(tr)
	var buf bytes.Buffer
	fmt := output.NewFormatter("json", &buf)

	n, err := runner.Run(strings.NewReader(""), parser.NewNginxParser(), f, fmt)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
	_ = time.Now() // keep import used
}
