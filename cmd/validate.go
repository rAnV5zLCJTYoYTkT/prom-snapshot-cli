package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	validateCheckOverlap      bool
	validateCheckCorruption   bool
	validateJSONOutput        bool
)

var validateCmd = &cobra.Command{
	Use:   "validate <snapshot-path>",
	Short: "Validate the structural integrity of a TSDB snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&validateCheckOverlap, "check-overlap", true, "check for overlapping blocks")
	validateCmd.Flags().BoolVar(&validateCheckCorruption, "check-corruption", false, "check for block corruption (slow)")
	validateCmd.Flags().BoolVar(&validateJSONOutput, "json", false, "output result as JSON")
}

func runValidate(cmd *cobra.Command, args []string) error {
	opts := snapshot.ValidationOptions{
		Path:            args[0],
		CheckOverlap:    validateCheckOverlap,
		CheckCorruption: validateCheckCorruption,
	}

	result, err := snapshot.ValidateSnapshot(opts)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if validateJSONOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Printf("Path:         %s\n", result.Path)
	fmt.Printf("Blocks found: %d\n", result.BlocksFound)
	fmt.Printf("Valid:        %v\n", result.Valid)
	if len(result.Issues) > 0 {
		fmt.Println("Issues:")
		for _, issue := range result.Issues {
			fmt.Printf("  - %s\n", issue)
		}
	}
	return nil
}
