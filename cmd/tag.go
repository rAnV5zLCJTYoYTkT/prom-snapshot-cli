package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	tagOutput     string
	tagLabelName  string
	tagLabelValue string
	tagMatchers   []string
	tagJSON       bool
)

var tagCmd = &cobra.Command{
	Use:   "tag <snapshot-path>",
	Short: "Add or overwrite a static label on matching series",
	Args:  cobra.ExactArgs(1),
	RunE:  runTag,
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.Flags().StringVarP(&tagOutput, "output", "o", "", "output directory for tagged snapshot (required)")
	tagCmd.Flags().StringVar(&tagLabelName, "label", "", "label name to add/overwrite (required)")
	tagCmd.Flags().StringVar(&tagLabelValue, "value", "", "label value to set (required)")
	tagCmd.Flags().StringArrayVar(&tagMatchers, "match", nil, "series selectors (repeatable)")
	tagCmd.Flags().BoolVar(&tagJSON, "json", false, "output result as JSON")
	_ = tagCmd.MarkFlagRequired("output")
	_ = tagCmd.MarkFlagRequired("label")
	_ = tagCmd.MarkFlagRequired("value")
}

func runTag(cmd *cobra.Command, args []string) error {
	matchers, err := snapshot.ParseMatchers(tagMatchers)
	if err != nil {
		return fmt.Errorf("invalid matcher: %w", err)
	}

	result, err := snapshot.TagSnapshot(snapshot.TagOptions{
		Path:       args[0],
		OutputDir:  tagOutput,
		LabelName:  tagLabelName,
		LabelValue: tagLabelValue,
		Matchers:   matchers,
	})
	if err != nil {
		return err
	}

	if tagJSON {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	fmt.Fprintf(os.Stdout, "Tagged:  %d / %d series\n", result.SeriesTagged, result.SeriesTotal)
	fmt.Fprintf(os.Stdout, "Output:  %s\n", result.OutputDir)
	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}
	return nil
}
