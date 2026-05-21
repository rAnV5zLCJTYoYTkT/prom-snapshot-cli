package snapshot

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
)

func TestDiffResult_Fields(t *testing.T) {
	d := DiffResult{
		Added:   []SeriesDiff{{Labels: map[string]string{"__name__": "up"}, OnlyIn: "B"}},
		Removed: []SeriesDiff{{Labels: map[string]string{"__name__": "down"}, OnlyIn: "A"}},
		Changed: []SeriesDiff{{Labels: map[string]string{"__name__": "go_gc"}, CountA: 5, CountB: 10}},
	}
	if len(d.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(d.Added))
	}
	if len(d.Removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(d.Removed))
	}
	if d.Changed[0].CountA != 5 || d.Changed[0].CountB != 10 {
		t.Fatalf("unexpected changed counts: %+v", d.Changed[0])
	}
}

func TestSeriesDiff_OnlyIn(t *testing.T) {
	sd := SeriesDiff{Labels: map[string]string{"job": "prometheus"}, OnlyIn: "A"}
	if sd.OnlyIn != "A" {
		t.Fatalf("expected OnlyIn=A, got %s", sd.OnlyIn)
	}
	if sd.Labels["job"] != "prometheus" {
		t.Fatalf("unexpected label value: %s", sd.Labels["job"])
	}
}

func TestIndexByFingerprint_Unique(t *testing.T) {
	series := []SeriesResult{
		{Labels: labels.FromStrings("__name__", "up", "job", "a"), SampleCount: 3},
		{Labels: labels.FromStrings("__name__", "up", "job", "b"), SampleCount: 5},
	}
	idx := indexByFingerprint(series)
	if len(idx) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(idx))
	}
}

func TestIndexByFingerprint_Duplicate(t *testing.T) {
	lbl := labels.FromStrings("__name__", "up")
	series := []SeriesResult{
		{Labels: lbl, SampleCount: 1},
		{Labels: lbl, SampleCount: 2},
	}
	idx := indexByFingerprint(series)
	// last write wins — still only one entry
	if len(idx) != 1 {
		t.Fatalf("expected 1 entry for duplicate fingerprint, got %d", len(idx))
	}
}

func TestDiffSnapshots_NilDB(t *testing.T) {
	_, err := DiffSnapshots(nil, nil, DiffOptions{
		Matchers: []*labels.Matcher{},
	})
	if err == nil {
		t.Fatal("expected error for nil db, got nil")
	}
}
