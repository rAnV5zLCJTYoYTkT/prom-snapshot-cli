package snapshot

import (
	"strings"
	"testing"
	"time"
)

func TestValidationResult_Fields(t *testing.T) {
	r := ValidationResult{
		Path:        "/tmp/snap",
		BlocksFound: 3,
		Issues:      []string{"overlap detected"},
		Valid:       false,
		CheckedAt:   time.Now(),
	}
	if r.Path != "/tmp/snap" {
		t.Errorf("expected path /tmp/snap, got %s", r.Path)
	}
	if r.BlocksFound != 3 {
		t.Errorf("expected 3 blocks, got %d", r.BlocksFound)
	}
	if r.Valid {
		t.Error("expected valid=false")
	}
	if len(r.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(r.Issues))
	}
}

func TestValidationResult_NoIssues(t *testing.T) {
	r := ValidationResult{
		Path:        "/tmp/snap",
		BlocksFound: 1,
		Valid:       true,
	}
	if !r.Valid {
		t.Error("expected valid=true")
	}
	if len(r.Issues) != 0 {
		t.Errorf("expected no issues, got %v", r.Issues)
	}
}

func TestValidateSnapshot_EmptyPath(t *testing.T) {
	_, err := ValidateSnapshot(ValidationOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
	if !strings.Contains(err.Error(), "path must not be empty") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSnapshot_InvalidPath(t *testing.T) {
	_, err := ValidateSnapshot(ValidationOptions{
		Path:         "/nonexistent/path/to/tsdb",
		CheckOverlap: true,
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestValidationOptions_Defaults(t *testing.T) {
	opts := ValidationOptions{
		Path: "/tmp/snap",
	}
	if opts.CheckOverlap {
		t.Error("CheckOverlap should default to false")
	}
	if opts.CheckCorruption {
		t.Error("CheckCorruption should default to false")
	}
}
