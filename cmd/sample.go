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
	sampleMatchers []string
	sampleStart    string
	sampleEnd      string
)

var sampleCmd = &cobra.Command{
	Use:   "sample <snapshot-path>",
	Short: "Fetch raw samples from a Prometheus TSDB snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runSample,
}

func init() {
	rootCmd.AddCommand(sampleCmd)
	sampleCmd.Flags().StringArrayVarP(&sampleMatchers, "match", "m", nil, "Label matchers (repeatable)")
	sampleCmd.Flags().StringVar(&sampleStart, "start", "", "Start time (RFC3339 or Unix seconds)")
	sampleCmd.Flags().StringVar(&sampleEnd, "end", "", "End time (RFC3339 or Unix seconds)")
	_ = sampleCmd.MarkFlagRequired("start")
	_ = sampleCmd.MarkFlagRequired("end")
}

func runSample(cmd *cobra.Command, args []string) error {
	start, err := parseTimeFlag(sampleStart)
	if err != nil {
		return fmt.Errorf("invalid --start: %w", err)
	}
	end, err := parseTimeFlag(sampleEnd)
	if err != nil {
		return fmt.Errorf("invalid --end: %w", err)
	}

	results, err := snapshot.FetchSamples(snapshot.SampleOptions{
		Path:     args[0],
		Matchers: sampleMatchers,
		Start:    start,
		End:      end,
		Step:     time.Minute,
	})
	if err != nil {
		return fmt.Errorf("fetching samples: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return fmt.Errorf("encoding output: %w", err)
	}
	return nil
}
