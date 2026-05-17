package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	snapshotPath string
	verbose      bool
)

var rootCmd = &cobra.Command{
	Use:   "prom-snapshot-cli",
	Short: "Query, export, and replay Prometheus TSDB snapshots",
	Long: `prom-snapshot-cli is a command-line utility for working with
Prometheus TSDB snapshots. It allows you to query metrics, export data,
and replay snapshots for local debugging purposes.`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&snapshotPath,
		"snapshot", "s",
		"",
		"Path to the Prometheus TSDB snapshot directory (required)",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose", "v",
		false,
		"Enable verbose output",
	)
}
