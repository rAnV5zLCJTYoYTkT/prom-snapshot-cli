package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	seriesMatchers []string
	seriesLimit    int
	seriesJSON     bool
)

var seriesCmd = &cobra.Command{
	Use:   "series <snapshot-path>",
	Short: "List series metadata from a Prometheus TSDB snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runSeries,
}

func init() {
	seriesCmd.Flags().StringArrayVarP(&seriesMatchers, "match", "m", nil, "Label matchers (e.g. '__name__=~\"up.*\"')")
	seriesCmd.Flags().IntVarP(&seriesLimit, "limit", "l", 100, "Maximum number of series to return")
	seriesCmd.Flags().BoolVar(&seriesJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(seriesCmd)
}

func runSeries(cmd *cobra.Command, args []string) error {
	path := args[0]

	matchers, err := snapshot.ParseMatchers(seriesMatchers)
	if err != nil {
		return fmt.Errorf("invalid matchers: %w", err)
	}

	db, err := snapshot.NewReader(path)
	if err != nil {
		return fmt.Errorf("opening snapshot: %w", err)
	}
	defer db.Close()

	summary, err := snapshot.ListSeries(db, matchers, seriesLimit)
	if err != nil {
		return fmt.Errorf("listing series: %w", err)
	}

	if seriesJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(summary)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "LABELS\tSAMPLES\tMIN_TIME_MS\tMAX_TIME_MS\n")
	for _, s := range summary.Series {
		lblStr := ""
		for k, v := range s.Labels {
			if lblStr != "" {
				lblStr += ","
			}
			lblStr += k + "=" + v
		}
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\n", lblStr, s.NumSamples, s.MinTime, s.MaxTime)
	}
	w.Flush()
	fmt.Fprintf(os.Stdout, "\nTotal series: %d | Total samples: %d\n", summary.TotalSeries, summary.TotalSamples)
	return nil
}
