package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	cardinalityTopN   int
	cardinalityJSON   bool
	cardinalityFilter string
)

var cardinalityCmd = &cobra.Command{
	Use:   "cardinality <snapshot-path>",
	Short: "Report label cardinality for a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runCardinality,
}

func init() {
	rootCmd.AddCommand(cardinalityCmd)
	cardinalityCmd.Flags().IntVar(&cardinalityTopN, "top", 10, "Number of top labels to report")
	cardinalityCmd.Flags().BoolVar(&cardinalityJSON, "json", false, "Output as JSON")
	cardinalityCmd.Flags().StringVar(&cardinalityFilter, "match", "", "Label matcher to filter series")
}

func runCardinality(cmd *cobra.Command, args []string) error {
	snapshotPath := args[0]

	matchers, err := snapshot.ParseMatchers(cardinalityFilter)
	if err != nil {
		return fmt.Errorf("invalid matcher: %w", err)
	}

	opts := snapshot.CardinalityOptions{
		Path:     snapshotPath,
		TopN:     cardinalityTopN,
		Matchers: matchers,
	}

	report, err := snapshot.CardinalitySnapshot(opts)
	if err != nil {
		return fmt.Errorf("cardinality error: %w", err)
	}

	if cardinalityJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	fmt.Fprintf(os.Stdout, "Total series: %d\n", report.TotalSeries)
	fmt.Fprintf(os.Stdout, "Top %d labels by series count:\n", len(report.TopN))
	for i, entry := range report.TopN {
		fmt.Fprintf(os.Stdout, "  %2d. %-40s %d\n", i+1, entry.LabelName, entry.SeriesCount)
	}
	return nil
}
