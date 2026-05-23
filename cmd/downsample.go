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
	downsampleOutput     string
	downsampleResolution time.Duration
	downsampleStart      string
	downsampleEnd        string
	downsampleMatchers   []string
	downsampleJSON       bool
)

func init() {
	downsampleCmd := &cobra.Command{
		Use:   "downsample <snapshot-path>",
		Short: "Downsample TSDB snapshot to a lower resolution",
		Args:  cobra.ExactArgs(1),
		RunE:  runDownsample,
	}
	downsampleCmd.Flags().StringVarP(&downsampleOutput, "output", "o", "", "output path (required)")
	downsampleCmd.Flags().DurationVarP(&downsampleResolution, "resolution", "r", 5*time.Minute, "bucket resolution (e.g. 5m, 1h)")
	downsampleCmd.Flags().StringVar(&downsampleStart, "start", "", "start time (RFC3339)")
	downsampleCmd.Flags().StringVar(&downsampleEnd, "end", "", "end time (RFC3339)")
	downsampleCmd.Flags().StringArrayVarP(&downsampleMatchers, "match", "m", nil, "label matchers")
	downsampleCmd.Flags().BoolVar(&downsampleJSON, "json", false, "output result as JSON")
	_ = downsampleCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(downsampleCmd)
}

func runDownsample(cmd *cobra.Command, args []string) error {
	matchers, err := snapshot.ParseMatchers(downsampleMatchers)
	if err != nil {
		return fmt.Errorf("invalid matcher: %w", err)
	}

	opts := snapshot.DownsampleOptions{
		Path:       args[0],
		OutputPath: downsampleOutput,
		Resolution: downsampleResolution,
		Matchers:   matchers,
	}

	if downsampleStart != "" {
		t, err := parseTimeFlag(downsampleStart)
		if err != nil {
			return fmt.Errorf("invalid start: %w", err)
		}
		opts.Start = t
	}
	if downsampleEnd != "" {
		t, err := parseTimeFlag(downsampleEnd)
		if err != nil {
			return fmt.Errorf("invalid end: %w", err)
		}
		opts.End = t
	}

	result, err := snapshot.DownsampleSnapshot(opts)
	if err != nil {
		return err
	}

	if downsampleJSON {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	fmt.Printf("Series processed : %d\n", result.SeriesProcessed)
	fmt.Printf("Samples in       : %d\n", result.SamplesIn)
	fmt.Printf("Samples out      : %d\n", result.SamplesOut)
	fmt.Printf("Reduction        : %.2f%%\n", result.ReductionPct)
	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}
	return nil
}
