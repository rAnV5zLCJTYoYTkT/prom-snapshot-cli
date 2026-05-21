package snapshot

import (
	"testing"
	"time"
)

func TestBlockStats_Fields(t *testing.T) {
	minT := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxT := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	bs := BlockStats{
		ULID:       "01HQ1234ABCD",
		MinTime:    minT,
		MaxTime:    maxT,
		Duration:   maxT.Sub(minT),
		NumSeries:  100,
		NumSamples: 5000,
		NumChunks:  200,
		SizeBytes:  1024 * 1024,
	}

	if bs.ULID != "01HQ1234ABCD" {
		t.Errorf("unexpected ULID: %s", bs.ULID)
	}
	if bs.Duration != 24*time.Hour {
		t.Errorf("expected 24h duration, got %v", bs.Duration)
	}
	if bs.NumSeries != 100 {
		t.Errorf("expected 100 series, got %d", bs.NumSeries)
	}
}

func TestSnapshotStats_Aggregation(t *testing.T) {
	now := time.Now().UTC()
	earlier := now.Add(-48 * time.Hour)

	stats := &SnapshotStats{
		Blocks: []BlockStats{
			{ULID: "A", MinTime: earlier, MaxTime: earlier.Add(24 * time.Hour), NumSeries: 50, NumSamples: 1000, NumChunks: 100, SizeBytes: 512},
			{ULID: "B", MinTime: now.Add(-24 * time.Hour), MaxTime: now, NumSeries: 80, NumSamples: 2000, NumChunks: 150, SizeBytes: 768},
		},
		TotalSeries:  130,
		TotalSamples: 3000,
		TotalChunks:  250,
		TotalBytes:   1280,
		EarliestTime: earlier,
		LatestTime:   now,
	}

	if stats.TotalSeries != 130 {
		t.Errorf("expected 130 total series, got %d", stats.TotalSeries)
	}
	if stats.TotalSamples != 3000 {
		t.Errorf("expected 3000 total samples, got %d", stats.TotalSamples)
	}
	if len(stats.Blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(stats.Blocks))
	}
	if !stats.EarliestTime.Equal(earlier) {
		t.Errorf("unexpected earliest time: %v", stats.EarliestTime)
	}
	if !stats.LatestTime.Equal(now) {
		t.Errorf("unexpected latest time: %v", stats.LatestTime)
	}
}

func TestSnapshotStats_Empty(t *testing.T) {
	stats := &SnapshotStats{}
	if stats.TotalSeries != 0 {
		t.Errorf("expected 0 series for empty stats")
	}
	if len(stats.Blocks) != 0 {
		t.Errorf("expected no blocks for empty stats")
	}
}
