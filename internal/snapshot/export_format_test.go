package snapshot

import (
	"testing"
)

func TestExportFormat_Constants(t *testing.T) {
	tests := []struct {
		format   ExportFormat
		expected string
	}{
		{FormatJSON, "json"},
		{FormatCSV, "csv"},
	}
	for _, tt := range tests {
		if string(tt.format) != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, string(tt.format))
		}
	}
}

func TestLabelsToMap(t *testing.T) {
	se := SeriesExport{
		Labels: map[string]string{
			"__name__": "http_requests_total",
			"job":      "api-server",
			"instance": "localhost:9090",
		},
		Samples: []Sample{},
	}
	if se.Labels["__name__"] != "http_requests_total" {
		t.Errorf("unexpected label value: %s", se.Labels["__name__"])
	}
	if len(se.Labels) != 3 {
		t.Errorf("expected 3 labels, got %d", len(se.Labels))
	}
}

func TestSample_Fields(t *testing.T) {
	import_time := mustParseRFC3339("2024-01-01T00:00:00Z")
	s := Sample{
		Timestamp: import_time,
		Value:     3.14,
	}
	if s.Value != 3.14 {
		t.Errorf("expected 3.14, got %f", s.Value)
	}
	if s.Timestamp != import_time {
		t.Errorf("unexpected timestamp: %v", s.Timestamp)
	}
}

func TestExportOptions_DefaultFormat(t *testing.T) {
	opts := ExportOptions{
		Format: FormatJSON,
	}
	if opts.Format != FormatJSON {
		t.Errorf("expected FormatJSON, got %q", opts.Format)
	}
}
