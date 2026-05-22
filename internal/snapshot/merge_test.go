package snapshot

import (
	"testing"
)

func TestMergeOptions_TooFewSources(t *testing.T) {
	_, err := MergeSnapshots(MergeOptions{
		SourcePaths: []string{"/tmp/only-one"},
		OutputPath:  "/tmp/out",
	})
	if err == nil {
		t.Fatal("expected error for fewer than 2 sources")
	}
}

func TestMergeOptions_EmptyOutput(t *testing.T) {
	_, err := MergeSnapshots(MergeOptions{
		SourcePaths: []string{"/tmp/a", "/tmp/b"},
		OutputPath:  "",
	})
	if err == nil {
		t.Fatal("expected error for empty output path")
	}
}

func TestMergeOptions_InvalidSourcePath(t *testing.T) {
	_, err := MergeSnapshots(MergeOptions{
		SourcePaths: []string{"/nonexistent/path/a", "/nonexistent/path/b"},
		OutputPath:  "/tmp/merge-out",
	})
	if err == nil {
		t.Fatal("expected error for invalid source paths")
	}
}

func TestMergeResult_Fields(t *testing.T) {
	r := MergeResult{
		SourceCount:   3,
		SeriesMerged:  42,
		BlocksWritten: 1,
		Warnings:      []string{"dup series 123"},
	}
	if r.SourceCount != 3 {
		t.Errorf("expected SourceCount 3, got %d", r.SourceCount)
	}
	if r.SeriesMerged != 42 {
		t.Errorf("expected SeriesMerged 42, got %d", r.SeriesMerged)
	}
	if r.BlocksWritten != 1 {
		t.Errorf("expected BlocksWritten 1, got %d", r.BlocksWritten)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestMergeResult_NoWarnings(t *testing.T) {
	r := MergeResult{
		SourceCount:   2,
		SeriesMerged:  10,
		BlocksWritten: 1,
	}
	if r.Warnings != nil {
		t.Errorf("expected nil warnings, got %v", r.Warnings)
	}
}
