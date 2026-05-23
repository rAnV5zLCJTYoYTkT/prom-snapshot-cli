package snapshot

import (
	"testing"
	"time"
)

func TestResampleOptions_MissingPath(t *testing.T) {
	_, err := ResampleSnapshot(ResampleOptions{
		Output:   "/tmp/out",
		Interval: 15 * time.Second,
	})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestResampleOptions_MissingOutput(t *testing.T) {
	_, err := ResampleSnapshot(ResampleOptions{
		Path:     "/tmp/snap",
		Interval: 15 * time.Second,
	})
	if err == nil || err.Error() != "output path is required" {
		t.Fatalf("expected 'output path is required', got %v", err)
	}
}

func TestResampleOptions_ZeroInterval(t *testing.T) {
	_, err := ResampleSnapshot(ResampleOptions{
		Path:     "/tmp/snap",
		Output:   "/tmp/out",
		Interval: 0,
	})
	if err == nil || err.Error() != "interval must be greater than zero" {
		t.Fatalf("expected interval error, got %v", err)
	}
}

func TestResampleOptions_EndBeforeStart(t *testing.T) {
	now := time.Now()
	_, err := ResampleSnapshot(ResampleOptions{
		Path:     "/tmp/snap",
		Output:   "/tmp/out",
		Interval: 30 * time.Second,
		Start:    now,
		End:      now.Add(-1 * time.Hour),
	})
	if err == nil || err.Error() != "end time must not be before start time" {
		t.Fatalf("expected end-before-start error, got %v", err)
	}
}

func TestResampleOptions_InvalidPath(t *testing.T) {
	_, err := ResampleSnapshot(ResampleOptions{
		Path:     "/nonexistent/path/snapshot",
		Output:   "/tmp/out",
		Interval: 15 * time.Second,
	})
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}

func TestResampleResult_Fields(t *testing.T) {
	r := ResampleResult{
		SeriesProcessed: 5,
		SamplesIn:       1000,
		SamplesOut:      200,
		Warnings:        []string{"warn1"},
	}
	if r.SeriesProcessed != 5 {
		t.Errorf("expected SeriesProcessed=5, got %d", r.SeriesProcessed)
	}
	if r.SamplesIn != 1000 {
		t.Errorf("expected SamplesIn=1000, got %d", r.SamplesIn)
	}
	if r.SamplesOut != 200 {
		t.Errorf("expected SamplesOut=200, got %d", r.SamplesOut)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestResampleResult_NoWarnings(t *testing.T) {
	r := ResampleResult{
		SeriesProcessed: 2,
		SamplesIn:       400,
		SamplesOut:      80,
	}
	if r.Warnings != nil {
		t.Errorf("expected nil warnings, got %v", r.Warnings)
	}
}
