package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInspectCmd_Registered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "inspect <snapshot-path>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("inspect command not registered")
	}
}

func TestInspectCmd_MissingArg(t *testing.T) {
	cmd := &cobra.Command{Use: "inspect", Args: cobra.ExactArgs(1), RunE: runInspect}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "")
	cmd.Flags().BoolVarP(&inspectVerbose, "verbose", "v", false, "")
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing argument")
	}
}

func TestInspectCmd_InvalidPath(t *testing.T) {
	cmd := &cobra.Command{Use: "inspect", Args: cobra.ExactArgs(1), RunE: runInspect}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "")
	cmd.Flags().BoolVarP(&inspectVerbose, "verbose", "v", false, "")
	cmd.SetArgs([]string{"/no/such/path"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestInspectCmd_JSONFlag(t *testing.T) {
	var found *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "inspect <snapshot-path>" {
			found = c
			break
		}
	}
	if found == nil {
		t.Fatal("inspect command not found")
	}
	if f := found.Flags().Lookup("json"); f == nil {
		t.Error("--json flag not registered on inspect command")
	}
}
