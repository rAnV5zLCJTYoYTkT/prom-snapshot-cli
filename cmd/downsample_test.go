package cmd

import (
	"testing"
)

func TestDownsampleCmd_Registered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "downsample <snapshot-path>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("downsample command not registered")
	}
}

func TestDownsampleCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"downsample", "--output", "/tmp/out"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing snapshot path")
	}
}

func TestDownsampleCmd_MissingOutputFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"downsample", "/tmp/snap"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing --output flag")
	}
}

func TestDownsampleCmd_InvalidMatcher(t *testing.T) {
	rootCmd.SetArgs([]string{
		"downsample", "/tmp/snap",
		"--output", "/tmp/out",
		"--match", "!!invalid",
	})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for invalid matcher")
	}
}

func TestDownsampleCmd_InvalidStart(t *testing.T) {
	rootCmd.SetArgs([]string{
		"downsample", "/tmp/snap",
		"--output", "/tmp/out",
		"--start", "not-a-time",
	})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for invalid start time")
	}
}

func TestDownsampleCmd_InvalidPath(t *testing.T) {
	rootCmd.SetArgs([]string{
		"downsample", "/nonexistent/snapshot",
		"--output", "/tmp/out",
		"--resolution", "1m",
	})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for invalid snapshot path")
	}
}
