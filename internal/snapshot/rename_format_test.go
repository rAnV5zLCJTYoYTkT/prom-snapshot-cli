package snapshot

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRenameResult_JSONRoundtrip(t *testing.T) {
	orig := RenameResult{
		SeriesModified: 7,
		OutputPath:     "/tmp/snap",
		Warnings:       []string{"w1", "w2"},
	}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got RenameResult
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SeriesModified != orig.SeriesModified {
		t.Errorf("SeriesModified mismatch: got %d", got.SeriesModified)
	}
	if got.OutputPath != orig.OutputPath {
		t.Errorf("OutputPath mismatch: got %s", got.OutputPath)
	}
	if len(got.Warnings) != 2 {
		t.Errorf("Warnings length mismatch: got %d", len(got.Warnings))
	}
}

func TestRenameResult_JSONKeys(t *testing.T) {
	r := RenameResult{SeriesModified: 3, OutputPath: "/out"}
	b, _ := json.Marshal(r)
	s := string(b)
	for _, key := range []string{"series_modified", "output_path"} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON key %q in %s", key, s)
		}
	}
}

func TestRenameResult_NoWarningsOmitted(t *testing.T) {
	r := RenameResult{SeriesModified: 0, OutputPath: "/out"}
	b, _ := json.Marshal(r)
	if strings.Contains(string(b), "warnings") {
		t.Errorf("expected warnings to be omitted when nil")
	}
}
