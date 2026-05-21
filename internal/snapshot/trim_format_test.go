package snapshot

import (
	"encoding/json"
	"testing"
)

func TestTrimResult_JSONRoundtrip(t *testing.T) {
	orig := TrimResult{
		BlocksInspected: 10,
		BlocksRetained:  7,
		BlocksDropped:   3,
		SeriesRetained:  42000,
		StartMs:         1700000000000,
		EndMs:           1700086400000,
	}

	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got TrimResult
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.BlocksInspected != orig.BlocksInspected {
		t.Errorf("BlocksInspected: want %d, got %d", orig.BlocksInspected, got.BlocksInspected)
	}
	if got.BlocksRetained != orig.BlocksRetained {
		t.Errorf("BlocksRetained: want %d, got %d", orig.BlocksRetained, got.BlocksRetained)
	}
	if got.BlocksDropped != orig.BlocksDropped {
		t.Errorf("BlocksDropped: want %d, got %d", orig.BlocksDropped, got.BlocksDropped)
	}
	if got.SeriesRetained != orig.SeriesRetained {
		t.Errorf("SeriesRetained: want %d, got %d", orig.SeriesRetained, got.SeriesRetained)
	}
	if got.StartMs != orig.StartMs {
		t.Errorf("StartMs: want %d, got %d", orig.StartMs, got.StartMs)
	}
	if got.EndMs != orig.EndMs {
		t.Errorf("EndMs: want %d, got %d", orig.EndMs, got.EndMs)
	}
}

func TestTrimResult_JSONKeys(t *testing.T) {
	r := TrimResult{BlocksInspected: 1, BlocksRetained: 1}
	b, _ := json.Marshal(r)
	s := string(b)

	for _, key := range []string{"blocks_inspected", "blocks_retained", "blocks_dropped", "series_retained", "start_ms", "end_ms"} {
		if !containsStr(s, key) {
			t.Errorf("JSON missing key %q in %s", key, s)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s[1:], sub) || s[:len(sub)] == sub)
}
