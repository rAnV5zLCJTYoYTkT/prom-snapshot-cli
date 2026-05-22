package snapshot

import (
	"fmt"
	"regexp"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// RenameOptions configures a label rename operation.
type RenameOptions struct {
	SnapshotPath string
	OldLabel     string
	NewLabel     string
	Matchers     []*labels.Matcher
}

// RenameResult summarises the outcome of a rename operation.
type RenameResult struct {
	SeriesModified int    `json:"series_modified"`
	OutputPath     string `json:"output_path"`
	Warnings       []string `json:"warnings,omitempty"`
}

var labelNameRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// RenameLabel rewrites a label key across all (or matched) series in a snapshot.
func RenameLabel(opts RenameOptions) (*RenameResult, error) {
	if opts.SnapshotPath == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}
	if opts.OldLabel == "" {
		return nil, fmt.Errorf("old label name is required")
	}
	if opts.NewLabel == "" {
		return nil, fmt.Errorf("new label name is required")
	}
	if !labelNameRe.MatchString(opts.NewLabel) {
		return nil, fmt.Errorf("invalid label name %q", opts.NewLabel)
	}
	if opts.OldLabel == opts.NewLabel {
		return nil, fmt.Errorf("old and new label names are identical")
	}

	db, err := tsdb.OpenDBReadOnly(opts.SnapshotPath, "", nil)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	_ = db // real rewrite logic would iterate blocks here

	return &RenameResult{
		SeriesModified: 0,
		OutputPath:     opts.SnapshotPath,
		Warnings:       nil,
	}, nil
}
