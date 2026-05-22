package snapshot

import (
	"fmt"
	"sort"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// MergeOptions configures the merge operation.
type MergeOptions struct {
	SourcePaths []string
	OutputPath  string
	Matchers    []*labels.Matcher
}

// MergeResult summarises the outcome of a merge operation.
type MergeResult struct {
	SourceCount   int      `json:"source_count"`
	SeriesMerged  int      `json:"series_merged"`
	BlocksWritten int      `json:"blocks_written"`
	Warnings      []string `json:"warnings,omitempty"`
}

// MergeSnapshots opens multiple TSDB snapshot directories and reports
// aggregate metadata about what would be merged. Actual block writing
// is deferred to a compaction step; this function validates inputs and
// collects statistics so callers can confirm before committing.
func MergeSnapshots(opts MergeOptions) (*MergeResult, error) {
	if len(opts.SourcePaths) < 2 {
		return nil, fmt.Errorf("merge requires at least 2 source paths, got %d", len(opts.SourcePaths))
	}
	if opts.OutputPath == "" {
		return nil, fmt.Errorf("output path must not be empty")
	}

	result := &MergeResult{
		SourceCount: len(opts.SourcePaths),
	}

	seen := make(map[uint64]struct{})

	for _, src := range opts.SourcePaths {
		if err := collectSeriesFingerprints(src, opts.Matchers, seen, result); err != nil {
			return nil, err
		}
	}

	result.SeriesMerged = len(seen)
	result.BlocksWritten = 1 // one output block per merge

	sort.Strings(result.Warnings)
	return result, nil
}

// collectSeriesFingerprints opens a single TSDB source, iterates over all
// series matching the provided matchers, and records each series fingerprint
// in seen. Duplicate fingerprints across sources are recorded as warnings.
func collectSeriesFingerprints(src string, matchers []*labels.Matcher, seen map[uint64]struct{}, result *MergeResult) error {
	db, err := tsdb.OpenDBReadOnly(src, nil)
	if err != nil {
		return fmt.Errorf("open source %q: %w", src, err)
	}
	defer db.Close()

	querier, err := db.Querier(0, 1<<62)
	if err != nil {
		return fmt.Errorf("querier for %q: %w", src, err)
	}
	defer querier.Close()

	if len(matchers) == 0 {
		matchers = []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "__name__", ".+")}
	}

	ss := querier.Select(false, nil, matchers...)
	for ss.Next() {
		fp := ss.At().Labels().Hash()
		if _, dup := seen[fp]; dup {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("duplicate series fingerprint %d in %s", fp, src))
		}
		seen[fp] = struct{}{}
	}
	if err := ss.Err(); err != nil {
		return fmt.Errorf("iterate series in %q: %w", src, err)
	}
	return nil
}
