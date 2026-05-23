package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var inspectVerbose bool

func init() {
	inspectCmd := &cobra.Command{
		Use:   "inspect <snapshot-path>",
		Short: "Show detailed per-block information for a TSDB snapshot",
		Args:  cobra.ExactArgs(1),
		RunE:  runInspect,
	}
	inspectCmd.Flags().BoolVarP(&inspectVerbose, "verbose", "v", false, "Include extra compaction metadata")
	inspectCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	report, err := snapshot.InspectSnapshot(snapshot.InspectOptions{
		Path:    args[0],
		Verbose: inspectVerbose,
	})
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	fmt.Printf("Snapshot : %s\n", report.Path)
	fmt.Printf("Blocks   : %d\n\n", report.Total)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ULID\tMIN TIME\tMAX TIME\tDURATION\tSERIES\tSAMPLES\tCHUNKS\tLEVEL")
	for _, b := range report.Blocks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\n",
			b.ULID,
			b.MinTime.Format("2006-01-02T15:04:05Z"),
			b.MaxTime.Format("2006-01-02T15:04:05Z"),
			b.Duration,
			b.Series,
			b.Samples,
			b.Chunks,
			b.Compaction,
		)
	}
	return w.Flush()
}
