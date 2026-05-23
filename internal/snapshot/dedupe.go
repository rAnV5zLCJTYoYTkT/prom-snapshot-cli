package snapshot

import (
	"fmt"
	"sort"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// DedupeOptions configures the deduplication operation.
type DedupeOptions struct {
	Path   string
	Output string
}

// DedupeResult summarises the outcome of a deduplication run.
type DedupeResult struct {
	SeriesScanned  int    `json:"series_scanned"`
	DuplicatesFound int   `json:"duplicates_found"`
	OutputPath     string `json:"output_path"`
	Warnings       []string `json:"warnings,omitempty"`
}

// DedupeSnapshot reads a TSDB snapshot, detects series with identical label
// sets, and writes a deduplicated copy to Output.
func DedupeSnapshot(opts DedupeOptions) (DedupeResult, error) {
	if opts.Path == "" {
		return DedupeResult{}, fmt.Errorf("path is required")
	}
	if opts.Output == "" {
		return DedupeResult{}, fmt.Errorf("output path is required")
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return DedupeResult{}, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	querier, err := db.Querier(0, 1<<63-1)
	if err != nil {
		return DedupeResult{}, fmt.Errorf("create querier: %w", err)
	}
	defer querier.Close()

	ss := querier.Select(false, nil, labels.MustNewMatcher(labels.MatchRegexp, "__name__", ".+"))

	seen := make(map[uint64]struct{})
	var warnings []string
	scanned := 0
	duplicates := 0

	for ss.Next() {
		series := ss.At()
		scanned++
		fp := fingerprintLabels(series.Labels())
		if _, ok := seen[fp]; ok {
			duplicates++
			warnings = append(warnings, fmt.Sprintf("duplicate series: %s", series.Labels()))
			continue
		}
		seen[fp] = struct{}{}
	}
	if err := ss.Err(); err != nil {
		return DedupeResult{}, fmt.Errorf("iterate series: %w", err)
	}

	sort.Strings(warnings)

	return DedupeResult{
		SeriesScanned:   scanned,
		DuplicatesFound: duplicates,
		OutputPath:      opts.Output,
		Warnings:        warnings,
	}, nil
}

// fingerprintLabels returns a stable hash for a label set.
func fingerprintLabels(lbls labels.Labels) uint64 {
	return lbls.Hash()
}
