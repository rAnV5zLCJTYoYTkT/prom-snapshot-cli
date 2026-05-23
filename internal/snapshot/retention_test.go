package snapshot

import (
	"testing"
	"time"
)

func TestRetentionOptions_EmptyPath(t *testing.T) {
	_, err := EnforceRetention(RetentionOptions{
		MaxAge: 24 * time.Hour,
	})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestRetentionOptions_ZeroMaxAge(t *testing.T) {
	_, err := EnforceRetention(RetentionOptions{
		Path:   "/some/path",
		MaxAge: 0,
	})
	if err == nil || err.Error() != "max-age must be a positive duration" {
		t.Fatalf("expected max-age error, got %v", err)
	}
}

func TestRetentionOptions_NegativeMaxAge(t *testing.T) {
	_, err := EnforceRetention(RetentionOptions{
		Path:   "/some/path",
		MaxAge: -time.Hour,
	})
	if err == nil {
		t.Fatal("expected error for negative max-age")
	}
}

func TestRetentionOptions_InvalidPath(t *testing.T) {
	_, err := EnforceRetention(RetentionOptions{
		Path:   "/nonexistent/snapshot/path",
		MaxAge: 24 * time.Hour,
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestRetentionResult_Fields(t *testing.T) {
	r := RetentionResult{
		BlocksScanned: 5,
		BlocksRemoved: 2,
		RemovedULIDs:  []string{"aaa", "bbb"},
		DryRun:        true,
	}
	if r.BlocksScanned != 5 {
		t.Errorf("BlocksScanned: want 5, got %d", r.BlocksScanned)
	}
	if r.BlocksRemoved != 2 {
		t.Errorf("BlocksRemoved: want 2, got %d", r.BlocksRemoved)
	}
	if len(r.RemovedULIDs) != 2 {
		t.Errorf("RemovedULIDs length: want 2, got %d", len(r.RemovedULIDs))
	}
	if !r.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestRetentionResult_NoRemovals(t *testing.T) {
	r := RetentionResult{
		BlocksScanned: 3,
		BlocksRemoved: 0,
	}
	if r.RemovedULIDs != nil {
		t.Error("expected nil RemovedULIDs when nothing removed")
	}
}
