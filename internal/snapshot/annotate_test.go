package snapshot

import (
	"testing"
)

func TestAnnotateOptions_MissingPath(t *testing.T) {
	_, err := AnnotateSnapshot(AnnotateOptions{
		OutputPath: "/tmp/out",
		LabelName:  "env",
		LabelValue: "prod",
	})
	if err == nil || err.Error() != "path is required" {
		t.Fatalf("expected 'path is required', got %v", err)
	}
}

func TestAnnotateOptions_MissingOutput(t *testing.T) {
	_, err := AnnotateSnapshot(AnnotateOptions{
		Path:       "/tmp/snap",
		LabelName:  "env",
		LabelValue: "prod",
	})
	if err == nil || err.Error() != "output path is required" {
		t.Fatalf("expected 'output path is required', got %v", err)
	}
}

func TestAnnotateOptions_MissingLabelName(t *testing.T) {
	_, err := AnnotateSnapshot(AnnotateOptions{
		Path:       "/tmp/snap",
		OutputPath: "/tmp/out",
		LabelValue: "prod",
	})
	if err == nil || err.Error() != "label name is required" {
		t.Fatalf("expected 'label name is required', got %v", err)
	}
}

func TestAnnotateOptions_InvalidLabelName(t *testing.T) {
	_, err := AnnotateSnapshot(AnnotateOptions{
		Path:       "/tmp/snap",
		OutputPath: "/tmp/out",
		LabelName:  "123invalid",
		LabelValue: "prod",
	})
	if err == nil {
		t.Fatal("expected error for invalid label name")
	}
}

func TestAnnotateOptions_MissingLabelValue(t *testing.T) {
	_, err := AnnotateSnapshot(AnnotateOptions{
		Path:       "/tmp/snap",
		OutputPath: "/tmp/out",
		LabelName:  "env",
	})
	if err == nil || err.Error() != "label value is required" {
		t.Fatalf("expected 'label value is required', got %v", err)
	}
}

func TestAnnotateOptions_InvalidPath(t *testing.T) {
	_, err := AnnotateSnapshot(AnnotateOptions{
		Path:       "/nonexistent/path",
		OutputPath: "/tmp/out",
		LabelName:  "env",
		LabelValue: "prod",
	})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestAnnotateResult_Fields(t *testing.T) {
	r := AnnotateResult{
		SeriesAnnotated: 10,
		SeriesSkipped:   2,
		Warnings:        []string{"w1"},
	}
	if r.SeriesAnnotated != 10 {
		t.Errorf("expected 10, got %d", r.SeriesAnnotated)
	}
	if r.SeriesSkipped != 2 {
		t.Errorf("expected 2, got %d", r.SeriesSkipped)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}
