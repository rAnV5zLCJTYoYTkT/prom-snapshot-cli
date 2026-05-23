package snapshot

import (
	"testing"
)

func TestTagOptions_MissingPath(t *testing.T) {
	_, err := TagSnapshot(TagOptions{
		OutputDir:  "/tmp/out",
		LabelName:  "env",
		LabelValue: "prod",
	})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestTagOptions_MissingOutput(t *testing.T) {
	_, err := TagSnapshot(TagOptions{
		Path:       "/tmp/snap",
		LabelName:  "env",
		LabelValue: "prod",
	})
	if err == nil || err.Error() != "output directory is required" {
		t.Fatalf("expected output error, got %v", err)
	}
}

func TestTagOptions_MissingLabelName(t *testing.T) {
	_, err := TagSnapshot(TagOptions{
		Path:       "/tmp/snap",
		OutputDir:  "/tmp/out",
		LabelValue: "prod",
	})
	if err == nil || err.Error() != "label name is required" {
		t.Fatalf("expected label name error, got %v", err)
	}
}

func TestTagOptions_InvalidLabelName(t *testing.T) {
	_, err := TagSnapshot(TagOptions{
		Path:       "/tmp/snap",
		OutputDir:  "/tmp/out",
		LabelName:  "123-bad!",
		LabelValue: "prod",
	})
	if err == nil {
		t.Fatal("expected invalid label name error")
	}
}

func TestTagOptions_InvalidLabelValue(t *testing.T) {
	_, err := TagSnapshot(TagOptions{
		Path:       "/tmp/snap",
		OutputDir:  "/tmp/out",
		LabelName:  "env",
		LabelValue: "bad value!",
	})
	if err == nil {
		t.Fatal("expected invalid label value error")
	}
}

func TestTagOptions_InvalidPath(t *testing.T) {
	_, err := TagSnapshot(TagOptions{
		Path:       "/nonexistent/path/snap",
		OutputDir:  "/tmp/out",
		LabelName:  "env",
		LabelValue: "prod",
	})
	if err == nil {
		t.Fatal("expected open snapshot error for invalid path")
	}
}

func TestTagResult_Fields(t *testing.T) {
	r := TagResult{
		SeriesTagged: 5,
		SeriesTotal:  10,
		OutputDir:    "/tmp/out",
		Warnings:     []string{"warn1"},
	}
	if r.SeriesTagged != 5 {
		t.Errorf("expected SeriesTagged=5, got %d", r.SeriesTagged)
	}
	if r.SeriesTotal != 10 {
		t.Errorf("expected SeriesTotal=10, got %d", r.SeriesTotal)
	}
	if r.OutputDir != "/tmp/out" {
		t.Errorf("unexpected OutputDir: %s", r.OutputDir)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}
