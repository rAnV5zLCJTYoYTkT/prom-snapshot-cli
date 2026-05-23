package snapshot

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDownsampleResult_JSONRoundtrip(t *testing.T) {
	orig := DownsampleResult{
		SeriesProcessed: 5,
		SamplesIn:       1000,
		SamplesOut:      200,
		ReductionPct:    80.0,
		Warnings:        []string{"stale chunk"},
	}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got DownsampleResult
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SeriesProcessed != orig.SeriesProcessed {
		t.Errorf("SeriesProcessed mismatch")
	}
	if got.SamplesIn != orig.SamplesIn {
		t.Errorf("SamplesIn mismatch")
	}
	if got.SamplesOut != orig.SamplesOut {
		t.Errorf("SamplesOut mismatch")
	}
	if got.ReductionPct != orig.ReductionPct {
		t.Errorf("ReductionPct mismatch")
	}
}

func TestDownsampleResult_JSONKeys(t *testing.T) {
	r := DownsampleResult{
		SeriesProcessed: 1,
		SamplesIn:       100,
		SamplesOut:      20,
		ReductionPct:    80.0,
	}
	b, _ := json.Marshal(r)
	s := string(b)
	for _, key := range []string{"series_processed", "samples_in", "samples_out", "reduction_pct"} {
		if !strings.Contains(s, key) {
			t.Errorf("missing JSON key %q in %s", key, s)
		}
	}
}

func TestDownsampleResult_WarningsOmittedWhenEmpty(t *testing.T) {
	r := DownsampleResult{
		SeriesProcessed: 1,
		SamplesIn:       50,
		SamplesOut:      10,
		ReductionPct:    80.0,
	}
	b, _ := json.Marshal(r)
	if strings.Contains(string(b), "warnings") {
		t.Errorf("warnings key should be omitted when empty")
	}
}
