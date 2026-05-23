package snapshot

import (
	"testing"
)

func TestCardinalityOptions_EmptyPath(t *testing.T) {
	_, err := CardinalitySnapshot(CardinalityOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCardinalityOptions_InvalidPath(t *testing.T) {
	_, err := CardinalitySnapshot(CardinalityOptions{Path: "/nonexistent/path", TopN: 5})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestCardinalityOptions_DefaultTopN(t *testing.T) {
	opts := CardinalityOptions{Path: "/tmp", TopN: 0}
	// We only verify that a zero TopN is normalised; actual DB open will fail
	// but the guard runs before the DB open — so we just check the error is
	// about opening the DB, not about TopN.
	_, err := CardinalitySnapshot(opts)
	if err == nil {
		t.Fatal("expected error opening /tmp as tsdb")
	}
	// Error should mention "open snapshot", not a TopN validation error.
	if err.Error() == "top-n must be positive" {
		t.Fatalf("unexpected TopN validation error: %v", err)
	}
}

func TestCardinalityEntry_Fields(t *testing.T) {
	e := CardinalityEntry{LabelName: "job", ValueCount: 42}
	if e.LabelName != "job" {
		t.Errorf("unexpected LabelName: %s", e.LabelName)
	}
	if e.ValueCount != 42 {
		t.Errorf("unexpected ValueCount: %d", e.ValueCount)
	}
}

func TestCardinalityReport_Fields(t *testing.T) {
	r := CardinalityReport{
		TotalSeries: 100,
		TopN: []CardinalityEntry{
			{LabelName: "__name__", ValueCount: 100},
			{LabelName: "job", ValueCount: 5},
		},
	}
	if r.TotalSeries != 100 {
		t.Errorf("unexpected TotalSeries: %d", r.TotalSeries)
	}
	if len(r.TopN) != 2 {
		t.Errorf("expected 2 entries, got %d", len(r.TopN))
	}
	if r.TopN[0].LabelName != "__name__" {
		t.Errorf("expected __name__ first, got %s", r.TopN[0].LabelName)
	}
}
