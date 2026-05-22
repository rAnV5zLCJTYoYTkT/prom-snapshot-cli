package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <source1> <source2> [sourceN...] --output <dir>",
	Short: "Merge multiple Prometheus TSDB snapshots into one",
	Args:  cobra.MinimumNArgs(2),
	RunE:  runMerge,
}

var (
	mergeOutput  string
	mergeJSON    bool
)

func init() {
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "output directory for merged snapshot (required)")
	_ = mergeCmd.MarkFlagRequired("output")
	mergeCmd.Flags().BoolVar(&mergeJSON, "json", false, "output result as JSON")
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
	for _, src := range args {
		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("source path %q not accessible: %w", src, err)
		}
	}

	opts := snapshot.MergeOptions{
		Sources:    args,
		OutputPath: mergeOutput,
	}

	result, err := snapshot.MergeSnapshots(opts)
	if err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	if mergeJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Fprintf(os.Stdout, "Merge complete\n")
	fmt.Fprintf(os.Stdout, "  Series written:  %d\n", result.SeriesWritten)
	fmt.Fprintf(os.Stdout, "  Samples written: %d\n", result.SamplesWritten)
	fmt.Fprintf(os.Stdout, "  Output path:     %s\n", result.OutputPath)

	if result.HasWarnings() {
		fmt.Fprintf(os.Stdout, "  Warnings (%d):\n", len(result.Warnings))
		for _, w := range result.Warnings {
			fmt.Fprintf(os.Stdout, "    - %s\n", w)
		}
	}

	return nil
}
