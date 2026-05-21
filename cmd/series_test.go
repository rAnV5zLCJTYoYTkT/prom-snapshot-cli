package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSeriesCmd_Registered(t *testing.T) {
	var found bool
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "series <snapshot-path>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("series command not registered on root")
	}
}

func TestSeriesCmd_MissingArg(t *testing.T) {
	cmd := &cobra.Command{Use: "series", RunE: runSeries, Args: cobra.ExactArgs(1)}
	cmd.Flags().StringArrayP("match", "m", nil, "")
	cmd.Flags().Int("limit", 100, "")
	cmd.Flags().Bool("json", false, "")
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing snapshot path")
	}
}

func TestSeriesCmd_InvalidMatcher(t *testing.T) {
	cmd := &cobra.Command{Use: "series", RunE: runSeries, Args: cobra.ExactArgs(1)}
	cmd.Flags().StringArrayP("match", "m", nil, "")
	cmd.Flags().Int("limit", 100, "")
	cmd.Flags().Bool("json", false, "")
	cmd.SetArgs([]string{"--match", "!!invalid", "/nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid matcher")
	}
}

func TestSeriesCmd_InvalidPath(t *testing.T) {
	cmd := &cobra.Command{Use: "series", RunE: runSeries, Args: cobra.ExactArgs(1)}
	cmd.Flags().StringArrayP("match", "m", nil, "")
	cmd.Flags().Int("limit", 100, "")
	cmd.Flags().Bool("json", false, "")
	cmd.SetArgs([]string{"/nonexistent/path/to/snapshot"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid snapshot path")
	}
}

func TestSeriesCmd_JSONFlag(t *testing.T) {
	f := seriesCmd.Flag("json")
	if f == nil {
		t.Fatal("--json flag not defined on series command")
	}
	if f.Value.Type() != "bool" {
		t.Errorf("expected bool flag, got %s", f.Value.Type())
	}
}
