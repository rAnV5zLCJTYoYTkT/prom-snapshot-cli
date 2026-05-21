package snapshot

import (
	"testing"
	"time"
)

func TestBlockMeta_Fields(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	bm := BlockMeta{
		ULID:       "01HXYZ1234567890ABCDEFGHIJ",
		MinTime:    now,
		MaxTime:    now.Add(2 * time.Hour),
		Duration:   2 * time.Hour,
		NumSeries:  100,
		NumSamples: 5000,
		NumChunks:  200,
		SizeBytes:  1024 * 1024,
	}

	if bm.ULID == "" {
		t.Error("expected non-empty ULID")
	}
	if bm.NumSeries != 100 {
		t.Errorf("expected NumSeries=100, got %d", bm.NumSeries)
	}
	if bm.Duration != 2*time.Hour {
		t.Errorf("expected Duration=2h, got %v", bm.Duration)
	}
}

func TestCompactSummary_Aggregation(t *testing.T) {
	s := &CompactSummary{
		Blocks: []BlockMeta{
			{NumSeries: 10, NumSamples: 100, SizeBytes: 512},
			{NumSeries: 20, NumSamples: 200, SizeBytes: 1024},
		},
		TotalSeries:    30,
		TotalSamples:   300,
		TotalSizeBytes: 1536,
		BlockCount:     2,
	}

	if s.BlockCount != 2 {
		t.Errorf("expected BlockCount=2, got %d", s.BlockCount)
	}
	if s.TotalSeries != 30 {
		t.Errorf("expected TotalSeries=30, got %d", s.TotalSeries)
	}
	if s.TotalSamples != 300 {
		t.Errorf("expected TotalSamples=300, got %d", s.TotalSamples)
	}
	if s.TotalSizeBytes != 1536 {
		t.Errorf("expected TotalSizeBytes=1536, got %d", s.TotalSizeBytes)
	}
}

func TestCompactSummary_Empty(t *testing.T) {
	s := &CompactSummary{
		Blocks:     []BlockMeta{},
		BlockCount: 0,
	}

	if len(s.Blocks) != 0 {
		t.Errorf("expected empty blocks slice, got %d", len(s.Blocks))
	}
	if s.TotalSeries != 0 {
		t.Errorf("expected TotalSeries=0, got %d", s.TotalSeries)
	}
}

func TestInspectBlocks_InvalidPath(t *testing.T) {
	_, err := InspectBlocks("/nonexistent/path/to/snapshot")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
