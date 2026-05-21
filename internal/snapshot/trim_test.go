package snapshot

import (
	"testing"
	"time"
)

func TestTrimOptions_InvalidPath(t *testing.T) {
	_, err := TrimSnapshot("", TrimOptions{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestTrimOptions_EndBeforeStart(t *testing.T) {
	now := time.Now()
	_, err := TrimSnapshot("/tmp/fake", TrimOptions{
		Start: now,
		End:   now.Add(-time.Minute),
	})
	if err == nil {
		t.Fatal("expected error when end <= start")
	}
}

func TestTrimOptions_EqualStartEnd(t *testing.T) {
	now := time.Now()
	_, err := TrimSnapshot("/tmp/fake", TrimOptions{
		Start: now,
		End:   now,
	})
	if err == nil {
		t.Fatal("expected error when end == start")
	}
}

func TestTrimResult_Fields(t *testing.T) {
	r := TrimResult{
		BlocksInspected: 5,
		BlocksRetained:  3,
		BlocksDropped:   2,
		SeriesRetained:  1000,
		StartMs:         1000,
		EndMs:           2000,
	}
	if r.BlocksInspected != r.BlocksRetained+r.BlocksDropped {
		t.Errorf("inspected=%d != retained(%d)+dropped(%d)",
			r.BlocksInspected, r.BlocksRetained, r.BlocksDropped)
	}
	if r.StartMs >= r.EndMs {
		t.Error("start must be before end")
	}
}

func TestTrimSnapshot_InvalidPath(t *testing.T) {
	_, err := TrimSnapshot("/nonexistent/path/db", TrimOptions{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
	})
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}
