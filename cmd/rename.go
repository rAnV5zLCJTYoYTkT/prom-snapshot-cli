package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	renameOldLabel string
	renameNewLabel string
	renameMatchers string
	renameJSON     bool
)

var renameCmd = &cobra.Command{
	Use:   "rename <snapshot-path>",
	Short: "Rename a label key across series in a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runRename,
}

func init() {
	renameCmd.Flags().StringVar(&renameOldLabel, "old", "", "label name to rename (required)")
	renameCmd.Flags().StringVar(&renameNewLabel, "new", "", "replacement label name (required)")
	renameCmd.Flags().StringVar(&renameMatchers, "match", "", "optional series selector")
	renameCmd.Flags().BoolVar(&renameJSON, "json", false, "output result as JSON")
	_ = renameCmd.MarkFlagRequired("old")
	_ = renameCmd.MarkFlagRequired("new")
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	matchers, err := snapshot.ParseMatchers(renameMatchers)
	if err != nil {
		return fmt.Errorf("invalid matcher: %w", err)
	}

	result, err := snapshot.RenameLabel(snapshot.RenameOptions{
		SnapshotPath: args[0],
		OldLabel:     renameOldLabel,
		NewLabel:     renameNewLabel,
		Matchers:     matchers,
	})
	if err != nil {
		return err
	}

	if renameJSON {
		return json.NewEncoder(os.Stdout).Encode(result)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Series modified: %d\nOutput: %s\n",
		result.SeriesModified, result.OutputPath)
	for _, w := range result.Warnings {
		fmt.Fprintf(cmd.ErrOrStderr(), "warning: %s\n", w)
	}
	return nil
}
