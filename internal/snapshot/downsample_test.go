package snapshot

import (
	"testing"
	"time"
)

func TestDownsampleOptions_MissingPath(t *testing.T) {
	_, err := DownsampleSnapshot(DownsampleOptions{
		OutputPath: "/tmp/out",
		Resolution: 5 * time.Minute,
	})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestDownsampleOptions_MissingOutput(t *testing.T) {
	_, err := DownsampleSnapshot(DownsampleOptions{
		Path:       "/tmp/snap",
		Resolution: 5 * time.Minute,
	})
	if err == nil || err.Error() != "output path is required" {
		t.Fatalf("expected output error, got %v", err)
	}
}

func TestDownsampleOptions_ZeroResolution(t *testing.T) {
	_, err := DownsampleSnapshot(DownsampleOptions{
		Path:       "/tmp/snap",
		OutputPath: "/tmp/out",
	})
	if err == nil || err.Error() != "resolution must be positive" {
		t.Fatalf("expected resolution error, got %v", err)
	}
}

func TestDownsampleOptions_StartAfterEnd(t *testing.T) {
	start := time.Now()
	end := start.Add(-time.Hour)
	_, err := DownsampleSnapshot(DownsampleOptions{
		Path:       "/tmp/snap",
		OutputPath: "/tmp/out",
		Resolution: time.Minute,
		Start:      start,
		End:        end,
	})
	if err == nil || err.Error() != "start must be before end" {
		t.Fatalf("expected time range error, got %v", err)
	}
}

func TestDownsampleOptions_InvalidPath(t *testing.T) {
	_, err := DownsampleSnapshot(DownsampleOptions{
		Path:       "/nonexistent/path",
		OutputPath: "/tmp/out",
		Resolution: time.Minute,
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestDownsampleResult_Fields(t *testing.T) {
	r := DownsampleResult{
		SeriesProcessed: 3,
		SamplesIn:       300,
		SamplesOut:      60,
		ReductionPct:    80.0,
		Warnings:        []string{"warn1"},
	}
	if r.SeriesProcessed != 3 {
		t.Errorf("expected 3 series, got %d", r.SeriesProcessed)
	}
	if r.ReductionPct != 80.0 {
		t.Errorf("expected 80.0 reduction, got %f", r.ReductionPct)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestBucketAverage_Empty(t *testing.T) {
	out := bucketAverage(nil, time.Minute)
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d samples", len(out))
	}
}

func TestBucketAverage_Aggregates(t *testing.T) {
	samples := []Sample{
		{TimestampMs: 0, Value: 10},
		{TimestampMs: 30000, Value: 20},
		{TimestampMs: 60000, Value: 30},
	}
	out := bucketAverage(samples, time.Minute)
	// bucket 0 covers 0..59999 → avg of 10 and 20 = 15
	// bucket 60000 covers 60000..119999 → 30
	if len(out) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(out))
	}
}
