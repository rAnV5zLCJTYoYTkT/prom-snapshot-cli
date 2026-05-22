package snapshot

import (
	"testing"
)

func TestRenameOptions_MissingPath(t *testing.T) {
	_, err := RenameLabel(RenameOptions{OldLabel: "job", NewLabel: "service"})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected path error, got %v", err)
	}
}

func TestRenameOptions_MissingOld(t *testing.T) {
	_, err := RenameLabel(RenameOptions{SnapshotPath: "/tmp/snap", NewLabel: "service"})
	if err == nil {
		t.Fatal("expected error for missing old label")
	}
}

func TestRenameOptions_MissingNew(t *testing.T) {
	_, err := RenameLabel(RenameOptions{SnapshotPath: "/tmp/snap", OldLabel: "job"})
	if err == nil {
		t.Fatal("expected error for missing new label")
	}
}

func TestRenameOptions_InvalidNewName(t *testing.T) {
	_, err := RenameLabel(RenameOptions{
		SnapshotPath: "/tmp/snap",
		OldLabel:     "job",
		NewLabel:     "123invalid",
	})
	if err == nil {
		t.Fatal("expected error for invalid label name")
	}
}

func TestRenameOptions_SameName(t *testing.T) {
	_, err := RenameLabel(RenameOptions{
		SnapshotPath: "/tmp/snap",
		OldLabel:     "job",
		NewLabel:     "job",
	})
	if err == nil {
		t.Fatal("expected error when old and new names are identical")
	}
}

func TestRenameOptions_InvalidPath(t *testing.T) {
	_, err := RenameLabel(RenameOptions{
		SnapshotPath: "/nonexistent/path",
		OldLabel:     "job",
		NewLabel:     "service",
	})
	if err == nil {
		t.Fatal("expected error for invalid snapshot path")
	}
}

func TestRenameResult_Fields(t *testing.T) {
	r := RenameResult{
		SeriesModified: 42,
		OutputPath:     "/tmp/out",
		Warnings:       []string{"warn1"},
	}
	if r.SeriesModified != 42 {
		t.Errorf("unexpected SeriesModified: %d", r.SeriesModified)
	}
	if r.OutputPath != "/tmp/out" {
		t.Errorf("unexpected OutputPath: %s", r.OutputPath)
	}
	if len(r.Warnings) != 1 {
		t.Errorf("unexpected Warnings length: %d", len(r.Warnings))
	}
}
