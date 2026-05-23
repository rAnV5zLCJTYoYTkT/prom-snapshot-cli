package snapshot_test

import (
	"encoding/json"
	"testing"
)

func TestCardinalityReport_JSONRoundtrip(t *testing.T) {
	report := CardinalityReport{
		TotalSeries: 42,
		TopN: []CardinalityEntry{
			{LabelName: "job", SeriesCount: 20},
			{LabelName: "instance", SeriesCount: 15},
		},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got CardinalityReport
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.TotalSeries != report.TotalSeries {
		t.Errorf("TotalSeries: got %d, want %d", got.TotalSeries, report.TotalSeries)
	}
	if len(got.TopN) != len(report.TopN) {
		t.Fatalf("TopN length: got %d, want %d", len(got.TopN), len(report.TopN))
	}
	if got.TopN[0].LabelName != "job" {
		t.Errorf("TopN[0].LabelName: got %q, want %q", got.TopN[0].LabelName, "job")
	}
}

func TestCardinalityReport_JSONKeys(t *testing.T) {
	report := CardinalityReport{
		TotalSeries: 5,
		TopN: []CardinalityEntry{
			{LabelName: "env", SeriesCount: 5},
		},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	for _, key := range []string{"total_series", "top_n"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected key %q in JSON output", key)
		}
	}
}

func TestCardinalityReport_TopNOmittedWhenEmpty(t *testing.T) {
	report := CardinalityReport{
		TotalSeries: 0,
		TopN:        []CardinalityEntry{},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if topN, ok := raw["top_n"]; ok {
		slice, _ := topN.([]interface{})
		if len(slice) != 0 {
			t.Errorf("expected empty top_n, got %v", topN)
		}
	}
}
