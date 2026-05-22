package snapshot_test

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMergeResult_JSONRoundtrip(t *testing.T) {
	r := MergeResult{
		SeriesWritten: 42,
		SamplesWritten: 1000,
		OutputPath:    "/tmp/merged",
		Warnings:      []string{"duplicate series detected"},
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got MergeResult
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.SeriesWritten != r.SeriesWritten {
		t.Errorf("SeriesWritten: got %d, want %d", got.SeriesWritten, r.SeriesWritten)
	}
	if got.SamplesWritten != r.SamplesWritten {
		t.Errorf("SamplesWritten: got %d, want %d", got.SamplesWritten, r.SamplesWritten)
	}
	if got.OutputPath != r.OutputPath {
		t.Errorf("OutputPath: got %q, want %q", got.OutputPath, r.OutputPath)
	}
	if len(got.Warnings) != 1 || got.Warnings[0] != r.Warnings[0] {
		t.Errorf("Warnings mismatch: got %v, want %v", got.Warnings, r.Warnings)
	}
}

func TestMergeResult_JSONKeys(t *testing.T) {
	r := MergeResult{
		SeriesWritten:  5,
		SamplesWritten: 50,
		OutputPath:    "/out",
		Warnings:      nil,
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	s := string(data)
	for _, key := range []string{"series_written", "samples_written", "output_path", "warnings"} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON key %q in output: %s", key, s)
		}
	}
}

func TestMergeResult_EmptyWarnings(t *testing.T) {
	r := MergeResult{
		SeriesWritten:  1,
		SamplesWritten: 10,
		OutputPath:    "/out",
		Warnings:      []string{},
	}

	if r.HasWarnings() {
		t.Error("expected HasWarnings to return false for empty slice")
	}
}

func TestMergeResult_WithWarnings(t *testing.T) {
	r := MergeResult{
		Warnings: []string{"warn1", "warn2"},
	}

	if !r.HasWarnings() {
		t.Error("expected HasWarnings to return true")
	}
}
