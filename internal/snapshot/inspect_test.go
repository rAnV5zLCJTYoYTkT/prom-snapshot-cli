package snapshot

import (
	"testing"
	"time"
)

func TestInspectOptions_EmptyPath(t *testing.T) {
	_, err := InspectSnapshot(InspectOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestInspectOptions_InvalidPath(t *testing.T) {
	_, err := InspectSnapshot(InspectOptions{Path: "/nonexistent/path"})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestBlockInspection_Fields(t *testing.T) {
	minT := time.Now().Add(-time.Hour)
	maxT := time.Now()
	insp := BlockInspection{
		ULID:       "01HXZ1234567890ABCDEFGHIJ",
		MinTime:    minT,
		MaxTime:    maxT,
		Duration:   maxT.Sub(minT).String(),
		Series:     42,
		Samples:    1000,
		Chunks:     10,
		Compaction: 1,
	}
	if insp.ULID == "" {
		t.Error("expected non-empty ULID")
	}
	if insp.Series != 42 {
		t.Errorf("expected 42 series, got %d", insp.Series)
	}
	if insp.Samples != 1000 {
		t.Errorf("expected 1000 samples, got %d", insp.Samples)
	}
	if insp.Duration == "" {
		t.Error("expected non-empty duration")
	}
}

func TestInspectReport_Fields(t *testing.T) {
	r := InspectReport{
		Path:   "/tmp/snapshot",
		Blocks: []BlockInspection{},
		Total:  0,
	}
	if r.Path != "/tmp/snapshot" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Total != 0 {
		t.Errorf("expected 0 total, got %d", r.Total)
	}
	if r.Blocks == nil {
		t.Error("expected non-nil blocks slice")
	}
}

func TestInspectReport_BlockCount(t *testing.T) {
	r := InspectReport{
		Blocks: []BlockInspection{
			{ULID: "A", Series: 1},
			{ULID: "B", Series: 2},
		},
		Total: 2,
	}
	if len(r.Blocks) != r.Total {
		t.Errorf("blocks len %d != total %d", len(r.Blocks), r.Total)
	}
}
