package snapshot

import (
	"testing"
	"time"
)

func TestSliceOptions_MissingPath(t *testing.T) {
	_, err := SliceSnapshot(SliceOptions{
		Output: "/tmp/out",
		Start:  time.Now().Add(-time.Hour),
		End:    time.Now(),
	})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestSliceOptions_MissingOutput(t *testing.T) {
	_, err := SliceSnapshot(SliceOptions{
		Path:  "/tmp/snap",
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
	})
	if err == nil || err.Error() != "output path is required" {
		t.Fatalf("expected 'output path is required', got %v", err)
	}
}

func TestSliceOptions_StartAfterEnd(t *testing.T) {
	now := time.Now()
	_, err := SliceSnapshot(SliceOptions{
		Path:   "/tmp/snap",
		Output: "/tmp/out",
		Start:  now,
		End:    now.Add(-time.Minute),
	})
	if err == nil {
		t.Fatal("expected error for start after end")
	}
}

func TestSliceOptions_InvalidPath(t *testing.T) {
	_, err := SliceSnapshot(SliceOptions{
		Path:   "/nonexistent/path/snap",
		Output: "/tmp/out",
		Start:  time.Now().Add(-time.Hour),
		End:    time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestSliceResult_Fields(t *testing.T) {
	r := SliceResult{
		SeriesWritten: 3,
		SamplesKept:   150,
		Warnings:      []string{"warn1"},
	}
	if r.SeriesWritten != 3 {
		t.Errorf("expected SeriesWritten=3, got %d", r.SeriesWritten)
	}
	if r.SamplesKept != 150 {
		t.Errorf("expected SamplesKept=150, got %d", r.SamplesKept)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestSliceResult_NoWarnings(t *testing.T) {
	r := SliceResult{
		SeriesWritten: 1,
		SamplesKept:   10,
	}
	if r.Warnings != nil {
		t.Errorf("expected nil warnings, got %v", r.Warnings)
	}
}
