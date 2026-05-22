package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestRenameCmd_Registered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("rename command not registered")
	}
}

func TestRenameCmd_MissingArg(t *testing.T) {
	cmd := &cobra.Command{Use: "rename", RunE: runRename, Args: cobra.ExactArgs(1)}
	cmd.Flags().StringVar(&renameOldLabel, "old", "job", "")
	cmd.Flags().StringVar(&renameNewLabel, "new", "service", "")
	cmd.Flags().StringVar(&renameMatchers, "match", "", "")
	cmd.Flags().BoolVar(&renameJSON, "json", false, "")
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot path")
	}
}

func TestRenameCmd_InvalidMatcher(t *testing.T) {
	cmd := &cobra.Command{Use: "rename", RunE: runRename, Args: cobra.ExactArgs(1)}
	cmd.Flags().StringVar(&renameOldLabel, "old", "job", "")
	cmd.Flags().StringVar(&renameNewLabel, "new", "service", "")
	cmd.Flags().StringVar(&renameMatchers, "match", "[[invalid", "")
	cmd.Flags().BoolVar(&renameJSON, "json", false, "")
	cmd.SetArgs([]string{"/tmp/snap"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid matcher")
	}
}

func TestRenameCmd_InvalidPath(t *testing.T) {
	cmd := &cobra.Command{Use: "rename", RunE: runRename, Args: cobra.ExactArgs(1)}
	cmd.Flags().StringVar(&renameOldLabel, "old", "job", "")
	cmd.Flags().StringVar(&renameNewLabel, "new", "service", "")
	cmd.Flags().StringVar(&renameMatchers, "match", "", "")
	cmd.Flags().BoolVar(&renameJSON, "json", false, "")
	cmd.SetArgs([]string{"/nonexistent/snap"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid snapshot path")
	}
}
