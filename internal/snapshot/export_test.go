package snapshot

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestWriteJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	if err := writeJSON(&buf, []SeriesExport{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got []SeriesExport
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(got))
	}
}

func TestWriteJSON_WithData(t *testing.T) {
	exports := []SeriesExport{
		{
			Labels: map[string]string{"__name__": "up", "job": "test"},
			Samples: []Sample{
				{Timestamp: time.UnixMilli(1000).UTC(), Value: 1.0},
				{Timestamp: time.UnixMilli(2000).UTC(), Value: 0.0},
			},
		},
	}
	var buf bytes.Buffer
	if err := writeJSON(&buf, exports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got []SeriesExport
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 series, got %d", len(got))
	}
	if len(got[0].Samples) != 2 {
		t.Errorf("expected 2 samples, got %d", len(got[0].Samples))
	}
}

func TestWriteCSV_Header(t *testing.T) {
	var buf bytes.Buffer
	if err := writeCSV(&buf, []SeriesExport{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "labels,timestamp,value") {
		t.Errorf("CSV missing header, got: %s", buf.String())
	}
}

func TestWriteCSV_WithData(t *testing.T) {
	exports := []SeriesExport{
		{
			Labels: map[string]string{"__name__": "up"},
			Samples: []Sample{
				{Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Value: 42.5},
			},
		},
	}
	var buf bytes.Buffer
	if err := writeCSV(&buf, exports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + 1 data), got %d", len(lines))
	}
	if !strings.Contains(lines[1], "42.5") {
		t.Errorf("expected value 42.5 in CSV line: %s", lines[1])
	}
}

func TestLabelsMapToString(t *testing.T) {
	m := map[string]string{"job": "prometheus"}
	s := labelsMapToString(m)
	if !strings.Contains(s, "prometheus") {
		t.Errorf("unexpected output: %s", s)
	}
}
