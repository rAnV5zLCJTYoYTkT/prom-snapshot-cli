package cmd

import (
	"testing"
	"time"
)

func TestTrimCmd_Registered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "trim <snapshot-path>" {
			return
		}
	}
	t.Fatal("trim command not registered")
}

func TestTrimCmd_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"trim", "--start", "2024-01-01T00:00:00Z", "--end", "2024-01-02T00:00:00Z"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when snapshot path is missing")
	}
}

func TestTrimCmd_MissingStartFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"trim", "/tmp/snap", "--end", "2024-01-02T00:00:00Z"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when --start is missing")
	}
}

func TestTrimCmd_MissingEndFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"trim", "/tmp/snap", "--start", "2024-01-01T00:00:00Z"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when --end is missing")
	}
}

func TestParseTimeFlag_RFC3339(t *testing.T) {
	s := "2024-06-01T12:00:00Z"
	got, err := parseTimeFlag(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Year() != 2024 || got.Month() != time.June {
		t.Errorf("unexpected time: %v", got)
	}
}

func TestParseTimeFlag_Unix(t *testing.T) {
	got, err := parseTimeFlag("1717200000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.IsZero() {
		t.Error("expected non-zero time")
	}
}

func TestParseTimeFlag_Invalid(t *testing.T) {
	_, err := parseTimeFlag("not-a-time")
	if err == nil {
		t.Fatal("expected error for invalid time string")
	}
}
