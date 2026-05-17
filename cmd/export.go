package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	exportFormat  string
	exportStart   string
	exportEnd     string
	exportMatch   []string
	exportOutput  string
)

var exportCmd = &cobra.Command{
	Use:   "export <snapshot-path>",
	Short: "Export time-series data from a Prometheus TSDB snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Output format: json or csv")
	exportCmd.Flags().StringVar(&exportStart, "start", "", "Start time (RFC3339 or Unix timestamp)")
	exportCmd.Flags().StringVar(&exportEnd, "end", "", "End time (RFC3339 or Unix timestamp)")
	exportCmd.Flags().StringArrayVar(&exportMatch, "match", nil, "Label matchers (e.g. '__name__=up')")
	exportCmd.Flags().StringVar(&exportOutput, "output", "-", "Output file path (- for stdout)")
}

func runExport(cmd *cobra.Command, args []string) error {
	snapshotPath := args[0]

	r, err := snapshot.NewReader(snapshotPath)
	if err != nil {
		return fmt.Errorf("failed to open snapshot: %w", err)
	}
	defer r.Close()

	matchers, err := snapshot.ParseMatchers(exportMatch)
	if err != nil {
		return fmt.Errorf("invalid matchers: %w", err)
	}

	qOpts := snapshot.QueryOptions{
		Matchers: matchers,
		Start:    time.Now().Add(-1 * time.Hour),
		End:      time.Now(),
	}

	if exportStart != "" {
		if t, err := time.Parse(time.RFC3339, exportStart); err == nil {
			qOpts.Start = t
		}
	}
	if exportEnd != "" {
		if t, err := time.Parse(time.RFC3339, exportEnd); err == nil {
			qOpts.End = t
		}
	}

	out := cmd.OutOrStdout()
	if exportOutput != "-" {
		f, err := os.Create(exportOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	return snapshot.ExportSeries(r, snapshot.ExportOptions{
		QueryOptions: qOpts,
		Format:       snapshot.ExportFormat(exportFormat),
		Out:          out,
	})
}
