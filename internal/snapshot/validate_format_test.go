package snapshot

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestValidationResult_JSONRoundtrip(t *testing.T) {
	orig := ValidationResult{
		Path:        "/data/snap",
		BlocksFound: 4,
		Issues:      []string{"block A overlaps block B"},
		Valid:       false,
		CheckedAt:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded ValidationResult
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Path != orig.Path {
		t.Errorf("path mismatch: got %s", decoded.Path)
	}
	if decoded.BlocksFound != orig.BlocksFound {
		t.Errorf("blocks_found mismatch: got %d", decoded.BlocksFound)
	}
	if decoded.Valid != orig.Valid {
		t.Errorf("valid mismatch")
	}
	if len(decoded.Issues) != len(orig.Issues) || decoded.Issues[0] != orig.Issues[0] {
		t.Errorf("issues mismatch: got %v, want %v", decoded.Issues, orig.Issues)
	}
	if !decoded.CheckedAt.Equal(orig.CheckedAt) {
		t.Errorf("checked_at mismatch: got %v, want %v", decoded.CheckedAt, orig.CheckedAt)
	}
}

func TestValidationResult_JSONKeys(t *testing.T) {
	r := ValidationResult{
		Path:        "/snap",
		BlocksFound: 2,
		Valid:       true,
		CheckedAt:   time.Now(),
	}
	b, _ := json.Marshal(r)
	s := string(b)

	for _, key := range []string{"path", "blocks_found", "valid", "checked_at"} {
		if !strings.Contains(s, key) {
			t.Errorf("JSON missing key %q", key)
		}
	}
}

func TestValidationResult_IssuesOmittedWhenEmpty(t *testing.T) {
	r := ValidationResult{
		Path:        "/snap",
		BlocksFound: 1,
		Valid:       true,
		CheckedAt:   time.Now(),
	}
	b, _ := json.Marshal(r)
	if strings.Contains(string(b), "issues") {
		t.Error("issues key should be omitted when empty")
	}
}
