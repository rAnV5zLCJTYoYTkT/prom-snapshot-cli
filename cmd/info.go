package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/prom-snapshot-cli/internal/snapshot"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display metadata and statistics about a TSDB snapshot",
	Long:  `Reads the TSDB snapshot directory and prints block metadata, time ranges, and series counts.`,
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(cmd *cobra.Command, args []string) error {
	if snapshotPath == "" {
		return fmt.Errorf("--snapshot flag is required")
	}

	reader, err := snapshot.NewReader(snapshotPath)
	if err != nil {
		return fmt.Errorf("failed to open snapshot: %w", err)
	}
	defer reader.Close()

	info, err := reader.Info()
	if err != nil {
		return fmt.Errorf("failed to read snapshot info: %w", err)
	}

	fmt.Printf("Snapshot Path : %s\n", snapshotPath)
	fmt.Printf("Blocks        : %d\n", info.BlockCount)
	fmt.Printf("Series        : %d\n", info.SeriesCount)
	fmt.Printf("Min Time      : %s\n", info.MinTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Max Time      : %s\n", info.MaxTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Printf("Duration      : %s\n", info.MaxTime.Sub(info.MinTime))

	if verbose {
		fmt.Println("\nBlocks:")
		for _, b := range info.Blocks {
			fmt.Printf("  [%s] %s → %s (%d series)\n",
				b.ULID,
				b.MinTime.Format("2006-01-02 15:04:05"),
				b.MaxTime.Format("2006-01-02 15:04:05"),
				b.NumSeries,
			)
		}
	}

	return nil
}
