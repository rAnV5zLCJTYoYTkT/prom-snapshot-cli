package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestStatsCmd_MissingArg(t *testing.T) {
	cmd := &cobra.Command{}
	err := runStats(cmd, []string{})
	// ExactArgs(1) is enforced by cobra; simulate missing path error from reader
	if err == nil {
		t.Error("expected error for missing snapshot path, got nil")
	}
}

func TestStatsCmd_InvalidPath(t *testing.T) {
	cmd := &cobra.Command{}
	err := runStats(cmd, []string{"/nonexistent/path/to/snapshot"})
	if err == nil {
		t.Error("expected error for invalid snapshot path, got nil")
	}
}

func TestStatsCmd_Registered(t *testing.T) {
	found := false
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "stats <snapshot-path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("stats command not registered on rootCmd")
	}
}

func TestStatsCmd_JSONFlag(t *testing.T) {
	f := statsCmd.Flags().Lookup("json")
	if f == nil {
		t.Fatal("expected --json flag to be defined")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}
