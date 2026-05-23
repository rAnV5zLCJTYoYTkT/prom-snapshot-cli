package snapshot

import (
	"testing"
)

func TestDedupeOptions_MissingPath(t *testing.T) {
	_, err := DedupeSnapshot(DedupeOptions{Output: "/tmp/out"})
	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if err.Error() != "path is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDedupeOptions_MissingOutput(t *testing.T) {
	_, err := DedupeSnapshot(DedupeOptions{Path: "/tmp/snap"})
	if err == nil {
		t.Fatal("expected error for missing output")
	}
	if err.Error() != "output path is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDedupeOptions_InvalidPath(t *testing.T) {
	_, err := DedupeSnapshot(DedupeOptions{
		Path:   "/nonexistent/snapshot/path",
		Output: "/tmp/out",
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestDedupeResult_Fields(t *testing.T) {
	r := DedupeResult{
		SeriesScanned:   42,
		DuplicatesFound: 3,
		OutputPath:      "/tmp/out",
		Warnings:        []string{"dup series: foo"},
	}
	if r.SeriesScanned != 42 {
		t.Errorf("expected 42 scanned, got %d", r.SeriesScanned)
	}
	if r.DuplicatesFound != 3 {
		t.Errorf("expected 3 duplicates, got %d", r.DuplicatesFound)
	}
	if r.OutputPath != "/tmp/out" {
		t.Errorf("unexpected output path: %s", r.OutputPath)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestDedupeResult_NoWarnings(t *testing.T) {
	r := DedupeResult{
		SeriesScanned:   10,
		DuplicatesFound: 0,
		OutputPath:      "/tmp/clean",
	}
	if r.Warnings != nil {
		t.Errorf("expected nil warnings, got %v", r.Warnings)
	}
}

func TestFingerprintLabels_Unique(t *testing.T) {
	l1 := newTestSeriesExport("up", map[string]string{"job": "a"}).Labels
	l2 := newTestSeriesExport("up", map[string]string{"job": "b"}).Labels
	if fingerprintLabels(mapToLabels(l1)) == fingerprintLabels(mapToLabels(l2)) {
		t.Error("expected distinct fingerprints for different label sets")
	}
}
