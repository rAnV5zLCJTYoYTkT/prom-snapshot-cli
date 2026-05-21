package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var statsJSON bool

var statsCmd = &cobra.Command{
	Use:   "stats <snapshot-path>",
	Short: "Show block-level statistics for a TSDB snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runStats,
}

func init() {
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "Output statistics as JSON")
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	snapshotPath := args[0]

	reader, err := snapshot.NewReader(snapshotPath)
	if err != nil {
		return fmt.Errorf("opening snapshot: %w", err)
	}
	defer reader.Close()

	stats, err := snapshot.CollectStats(reader.DB)
	if err != nil {
		return fmt.Errorf("collecting stats: %w", err)
	}

	if statsJSON {
		return json.NewEncoder(os.Stdout).Encode(stats)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Snapshot: %s\n", snapshotPath)
	fmt.Fprintf(w, "Blocks:\t%d\n", len(stats.Blocks))
	fmt.Fprintf(w, "Total series:\t%d\n", stats.TotalSeries)
	fmt.Fprintf(w, "Total samples:\t%d\n", stats.TotalSamples)
	fmt.Fprintf(w, "Total chunks:\t%d\n", stats.TotalChunks)
	fmt.Fprintf(w, "Total size:\t%d bytes\n", stats.TotalBytes)
	fmt.Fprintf(w, "Earliest time:\t%s\n", stats.EarliestTime.Format("2006-01-02T15:04:05Z"))
	fmt.Fprintf(w, "Latest time:\t%s\n", stats.LatestTime.Format("2006-01-02T15:04:05Z"))
	w.Flush()

	if len(stats.Blocks) > 0 {
		fmt.Fprintln(os.Stdout, "\nBlock details:")
		w2 := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w2, "ULID\tMIN TIME\tMAX TIME\tSERIES\tSAMPLES\tSIZE")
		for _, b := range stats.Blocks {
			fmt.Fprintf(w2, "%s\t%s\t%s\t%d\t%d\t%d\n",
				b.ULID,
				b.MinTime.Format("2006-01-02T15:04:05Z"),
				b.MaxTime.Format("2006-01-02T15:04:05Z"),
				b.NumSeries,
				b.NumSamples,
				b.SizeBytes,
			)
		}
		w2.Flush()
	}
	return nil
}
