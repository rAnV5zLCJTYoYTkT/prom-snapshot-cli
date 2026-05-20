package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplayCmd_MissingArg(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"replay"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing snapshot-dir argument")
	}
}

func TestReplayCmd_InvalidMatcher(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"replay", "/tmp/snap", "--match", "!!invalid!!"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid matcher")
	}
	if !strings.Contains(err.Error(), "parse matchers") && !strings.Contains(err.Error(), "matcher") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestReplayCmd_InvalidStart(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"replay", "/tmp/snap", "--start", "not-a-date"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --start")
	}
}
