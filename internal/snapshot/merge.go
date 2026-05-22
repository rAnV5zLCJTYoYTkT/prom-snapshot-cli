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
	SourceCount int      `json:"source_count"`
	SeriesMerged int     `json:"series_merged"`
	BlocksWritten int    `json:"blocks_written"`
	Warnings     []string `json:"warnings,omitempty"`
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
		db, err := tsdb.OpenDBReadOnly(src, nil)
		if err != nil {
			return nil, fmt.Errorf("open source %q: %w", src, err)
		}

		querier, err := db.Querier(0, 1<<62)
		if err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("querier for %q: %w", src, err)
		}

		matchers := opts.Matchers
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
			_ = querier.Close()
			_ = db.Close()
			return nil, fmt.Errorf("iterate series in %q: %w", src, err)
		}
		_ = querier.Close()
		_ = db.Close()
	}

	result.SeriesMerged = len(seen)
	result.BlocksWritten = 1 // one output block per merge

	sort.Strings(result.Warnings)
	return result, nil
}
