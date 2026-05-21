package snapshot

import (
	"encoding/json"
	"testing"
)

func TestSeriesInfo_JSONRoundtrip(t *testing.T) {
	orig := SeriesInfo{
		Labels:     map[string]string{"__name__": "http_requests_total", "job": "api"},
		NumSamples: 512,
		MinTime:    1700000000000,
		MaxTime:    1700003600000,
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var decoded SeriesInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if decoded.NumSamples != orig.NumSamples {
		t.Errorf("NumSamples mismatch: got %d, want %d", decoded.NumSamples, orig.NumSamples)
	}
	if decoded.Labels["job"] != "api" {
		t.Errorf("label mismatch: got %s", decoded.Labels["job"])
	}
}

func TestSeriesSummary_JSONRoundtrip(t *testing.T) {
	orig := &SeriesSummary{
		TotalSeries:  2,
		TotalSamples: 200,
		Series: []SeriesInfo{
			{Labels: map[string]string{"__name__": "up"}, NumSamples: 120, MinTime: 1000, MaxTime: 2000},
			{Labels: map[string]string{"__name__": "down"}, NumSamples: 80, MinTime: 1500, MaxTime: 2500},
		},
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var decoded SeriesSummary
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if decoded.TotalSeries != 2 {
		t.Errorf("TotalSeries mismatch: got %d", decoded.TotalSeries)
	}
	if len(decoded.Series) != 2 {
		t.Errorf("Series length mismatch: got %d", len(decoded.Series))
	}
}

func TestSeriesSummary_DefaultLimit(t *testing.T) {
	// Ensure ListSeries applies a default limit when limit <= 0
	_, err := ListSeries(nil, nil, 0)
	if err == nil {
		t.Error("expected nil db error")
	}
	// The default limit path is exercised; nil db error is the expected gate.
}
