package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	labelName     string
	labelMatchers []string
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "List label names or values from a snapshot",
	RunE:  runLabels,
}

func init() {
	rootCmd.AddCommand(labelsCmd)
	labelsCmd.Flags().StringVar(&labelName, "name", "", "Label name to list values for (omit to list all label names)")
	labelsCmd.Flags().StringArrayVar(&labelMatchers, "match", nil, "Label matchers to filter series (e.g. 'job=\"prometheus\"')")
}

func runLabels(cmd *cobra.Command, args []string) error {
	if snapshotPath == "" {
		return fmt.Errorf("--snapshot path is required")
	}

	r, err := snapshot.NewReader(snapshotPath)
	if err != nil {
		return fmt.Errorf("opening snapshot: %w", err)
	}
	defer r.Close()

	matchers, err := snapshot.ParseMatchers(labelMatchers)
	if err != nil {
		return fmt.Errorf("parsing matchers: %w", err)
	}

	lsets, err := r.LabelSets(matchers)
	if err != nil {
		return fmt.Errorf("reading label sets: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	if labelName != "" {
		fmt.Fprintf(w, "VALUES FOR LABEL\t%s\n", labelName)
		fmt.Fprintln(w, "---")
		for _, v := range snapshot.LabelValues(lsets, labelName) {
			fmt.Fprintln(w, v)
		}
	} else {
		fmt.Fprintln(w, "LABEL NAMES")
		fmt.Fprintln(w, "---")
		for _, n := range snapshot.LabelNames(lsets) {
			fmt.Fprintln(w, n)
		}
	}
	return nil
}
