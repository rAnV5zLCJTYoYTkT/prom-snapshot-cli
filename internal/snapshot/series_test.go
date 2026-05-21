package snapshot

import (
	"testing"
)

func TestSeriesInfo_Fields(t *testing.T) {
	info := SeriesInfo{
		Labels:     map[string]string{"__name__": "up", "job": "test"},
		NumSamples: 42,
		MinTime:    1000,
		MaxTime:    9000,
	}
	if info.NumSamples != 42 {
		t.Errorf("expected 42 samples, got %d", info.NumSamples)
	}
	if info.Labels["job"] != "test" {
		t.Errorf("expected job=test, got %s", info.Labels["job"])
	}
	if info.MinTime != 1000 || info.MaxTime != 9000 {
		t.Errorf("unexpected time range: %d-%d", info.MinTime, info.MaxTime)
	}
}

func TestSeriesSummary_Aggregation(t *testing.T) {
	summary := &SeriesSummary{
		TotalSeries:  3,
		TotalSamples: 150,
		Series: []SeriesInfo{
			{Labels: map[string]string{"__name__": "a"}, NumSamples: 100},
			{Labels: map[string]string{"__name__": "b"}, NumSamples: 30},
			{Labels: map[string]string{"__name__": "c"}, NumSamples: 20},
		},
	}
	if summary.TotalSeries != 3 {
		t.Errorf("expected 3 series, got %d", summary.TotalSeries)
	}
	if summary.TotalSamples != 150 {
		t.Errorf("expected 150 samples, got %d", summary.TotalSamples)
	}
	if len(summary.Series) != 3 {
		t.Errorf("expected 3 series entries, got %d", len(summary.Series))
	}
}

func TestSeriesSummary_Empty(t *testing.T) {
	summary := &SeriesSummary{
		TotalSeries:  0,
		TotalSamples: 0,
		Series:       nil,
	}
	if summary.TotalSeries != 0 {
		t.Errorf("expected 0 series")
	}
	if len(summary.Series) != 0 {
		t.Errorf("expected empty series slice")
	}
}

func TestListSeries_NilDB(t *testing.T) {
	_, err := ListSeries(nil, nil, 10)
	if err == nil {
		t.Error("expected error for nil db, got nil")
	}
}
