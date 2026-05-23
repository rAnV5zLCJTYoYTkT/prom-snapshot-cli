package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate <snapshot-path>",
	Short: "Add a constant label to matching series in a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runAnnotate,
}

func init() {
	annotateCmd.Flags().StringP("output", "o", "", "output path for annotated snapshot (required)")
	annotateCmd.Flags().StringP("label-name", "n", "", "label name to add (required)")
	annotateCmd.Flags().StringP("label-value", "v", "", "label value to set (required)")
	annotateCmd.Flags().StringSliceP("match", "m", nil, "series selectors to filter (repeatable)")
	annotateCmd.Flags().String("start", "", "start time (RFC3339 or unix seconds)")
	annotateCmd.Flags().String("end", "", "end time (RFC3339 or unix seconds)")
	annotateCmd.Flags().Bool("json", false, "output result as JSON")
	_ = annotateCmd.MarkFlagRequired("output")
	_ = annotateCmd.MarkFlagRequired("label-name")
	_ = annotateCmd.MarkFlagRequired("label-value")
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	labelName, _ := cmd.Flags().GetString("label-name")
	labelValue, _ := cmd.Flags().GetString("label-value")
	matchStrs, _ := cmd.Flags().GetStringSlice("match")
	startStr, _ := cmd.Flags().GetString("start")
	endStr, _ := cmd.Flags().GetString("end")
	asJSON, _ := cmd.Flags().GetBool("json")

	matchers, err := snapshot.ParseMatchers(matchStrs)
	if err != nil {
		return fmt.Errorf("invalid matcher: %w", err)
	}

	var start, end time.Time
	if startStr != "" {
		start, err = parseTimeFlag(startStr)
		if err != nil {
			return fmt.Errorf("invalid --start: %w", err)
		}
	}
	if endStr != "" {
		end, err = parseTimeFlag(endStr)
		if err != nil {
			return fmt.Errorf("invalid --end: %w", err)
		}
	}

	result, err := snapshot.AnnotateSnapshot(snapshot.AnnotateOptions{
		Path:       args[0],
		OutputPath: output,
		LabelName:  labelName,
		LabelValue: labelValue,
		Matchers:   matchers,
		Start:      start,
		End:        end,
	})
	if err != nil {
		return err
	}

	if asJSON {
		return json.NewEncoder(os.Stdout).Encode(result)
	}
	fmt.Printf("annotated: %d series, skipped: %d\n", result.SeriesAnnotated, result.SeriesSkipped)
	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}
	return nil
}
