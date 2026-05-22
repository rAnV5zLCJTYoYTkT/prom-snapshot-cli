package snapshot

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFilterOptions_MissingPath(t *testing.T) {
	_, err := FilterSnapshot(FilterOptions{Output: "/tmp/out"})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestFilterOptions_MissingOutput(t *testing.T) {
	_, err := FilterSnapshot(FilterOptions{Path: "/tmp/snap"})
	if err == nil || err.Error() != "output path is required" {
		t.Fatalf("expected 'output path is required', got %v", err)
	}
}

func TestFilterOptions_StartAfterEnd(t *testing.T) {
	now := time.Now()
	_, err := FilterSnapshot(FilterOptions{
		Path:   "/tmp/snap",
		Output: "/tmp/out",
		Start:  now.Add(time.Hour),
		End:    now,
	})
	if err == nil {
		t.Fatal("expected error for start after end")
	}
}

func TestFilterOptions_InvalidPath(t *testing.T) {
	_, err := FilterSnapshot(FilterOptions{
		Path:   "/nonexistent/path",
		Output: "/tmp/out",
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestFilterResult_Fields(t *testing.T) {
	r := FilterResult{
		SeriesKept:    10,
		SeriesDropped: 2,
		SamplesKept:   500,
		OutputPath:    "/tmp/out",
	}
	if r.SeriesKept != 10 {
		t.Errorf("expected SeriesKept=10, got %d", r.SeriesKept)
	}
	if r.SeriesDropped != 2 {
		t.Errorf("expected SeriesDropped=2, got %d", r.SeriesDropped)
	}
	if r.SamplesKept != 500 {
		t.Errorf("expected SamplesKept=500, got %d", r.SamplesKept)
	}
	if r.OutputPath != "/tmp/out" {
		t.Errorf("unexpected OutputPath: %s", r.OutputPath)
	}
}

func TestFilterResult_JSONRoundtrip(t *testing.T) {
	orig := FilterResult{
		SeriesKept:    3,
		SeriesDropped: 1,
		SamplesKept:   42,
		OutputPath:    "/data/filtered",
	}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got FilterResult
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got != orig {
		t.Errorf("roundtrip mismatch: got %+v, want %+v", got, orig)
	}
}

func TestFilterResult_JSONKeys(t *testing.T) {
	r := FilterResult{SeriesKept: 1, SeriesDropped: 0, SamplesKept: 10, OutputPath: "/out"}
	b, _ := json.Marshal(r)
	s := string(b)
	for _, key := range []string{"series_kept", "series_dropped", "samples_kept", "output_path"} {
		if !containsStr(s, key) {
			t.Errorf("missing JSON key %q in %s", key, s)
		}
	}
}
