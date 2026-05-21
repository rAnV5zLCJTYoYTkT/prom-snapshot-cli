package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	trimStart  string
	trimEnd    string
	trimJSON   bool
)

var trimCmd = &cobra.Command{
	Use:   "trim <snapshot-path>",
	Short: "Report which TSDB blocks fall within a time window",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrim,
}

func init() {
	trimCmd.Flags().StringVar(&trimStart, "start", "", "Start time (RFC3339 or Unix seconds) (required)")
	trimCmd.Flags().StringVar(&trimEnd, "end", "", "End time (RFC3339 or Unix seconds) (required)")
	trimCmd.Flags().BoolVar(&trimJSON, "json", false, "Output result as JSON")
	_ = trimCmd.MarkFlagRequired("start")
	_ = trimCmd.MarkFlagRequired("end")
	rootCmd.AddCommand(trimCmd)
}

func runTrim(cmd *cobra.Command, args []string) error {
	snapshotPath := args[0]

	start, err := parseTimeFlag(trimStart)
	if err != nil {
		return fmt.Errorf("invalid --start: %w", err)
	}
	end, err := parseTimeFlag(trimEnd)
	if err != nil {
		return fmt.Errorf("invalid --end: %w", err)
	}

	result, err := snapshot.TrimSnapshot(snapshotPath, snapshot.TrimOptions{
		Start: start,
		End:   end,
	})
	if err != nil {
		return fmt.Errorf("trim: %w", err)
	}

	if trimJSON {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	fmt.Fprintf(os.Stdout, "Blocks inspected : %d\n", result.BlocksInspected)
	fmt.Fprintf(os.Stdout, "Blocks retained  : %d\n", result.BlocksRetained)
	fmt.Fprintf(os.Stdout, "Blocks dropped   : %d\n", result.BlocksDropped)
	fmt.Fprintf(os.Stdout, "Series retained  : %d\n", result.SeriesRetained)
	return nil
}

// parseTimeFlag accepts RFC3339 or Unix-second strings.
func parseTimeFlag(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	var unix int64
	if _, err := fmt.Sscanf(s, "%d", &unix); err == nil {
		return time.Unix(unix, 0).UTC(), nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as RFC3339 or Unix seconds", s)
}
