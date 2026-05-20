package timerange_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timerange"
)

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNew_ValidRange(t *testing.T) {
	start := mustParse("2024-01-01T00:00:00Z")
	end := mustParse("2024-01-02T00:00:00Z")
	tr, err := timerange.New(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.Start.Equal(start) || !tr.End.Equal(end) {
		t.Errorf("expected [%v, %v], got %v", start, end, tr)
	}
}

func TestNew_InvalidRange(t *testing.T) {
	start := mustParse("2024-01-02T00:00:00Z")
	end := mustParse("2024-01-01T00:00:00Z")
	_, err := timerange.New(start, end)
	if err == nil {
		t.Fatal("expected error for start after end, got nil")
	}
}

func TestParse_Valid(t *testing.T) {
	tr, err := timerange.Parse("2024-03-01T10:00:00Z", "2024-03-01T12:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.IsZero() {
		t.Error("expected non-zero TimeRange")
	}
}

func TestParse_InvalidStart(t *testing.T) {
	_, err := timerange.Parse("not-a-time", "2024-03-01T12:00:00Z")
	if err == nil {
		t.Fatal("expected error for invalid start")
	}
}

func TestContains(t *testing.T) {
	tr, _ := timerange.Parse("2024-01-01T00:00:00Z", "2024-01-01T23:59:59Z")

	cases := []struct {
		ts       string
		inside   bool
	}{
		{"2024-01-01T00:00:00Z", true},
		{"2024-01-01T12:00:00Z", true},
		{"2024-01-01T23:59:59Z", true},
		{"2023-12-31T23:59:59Z", false},
		{"2024-01-02T00:00:00Z", false},
	}

	for _, c := range cases {
		t.Run(c.ts, func(t *testing.T) {
			got := tr.Contains(mustParse(c.ts))
			if got != c.inside {
				t.Errorf("Contains(%s) = %v, want %v", c.ts, got, c.inside)
			}
		})
	}
}
