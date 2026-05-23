package snapshot

import (
	"fmt"
	"sort"

	"github.com/prometheus/prometheus/tsdb"
)

// CardinalityEntry holds the label name and the count of unique values.
type CardinalityEntry struct {
	LabelName  string `json:"label_name"`
	ValueCount int    `json:"value_count"`
}

// CardinalityReport summarises label cardinality across a snapshot.
type CardinalityReport struct {
	TotalSeries int                `json:"total_series"`
	TopN        []CardinalityEntry `json:"top_n"`
}

// CardinalityOptions controls the cardinality analysis.
type CardinalityOptions struct {
	Path  string
	TopN  int
}

// CardinalitySnapshot opens the TSDB at path and returns a cardinality report
// listing the top-N label names by number of unique values.
func CardinalitySnapshot(opts CardinalityOptions) (*CardinalityReport, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("snapshot path must not be empty")
	}
	if opts.TopN <= 0 {
		opts.TopN = 10
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	querier, err := db.Querier(0, 1<<63-1)
	if err != nil {
		return nil, fmt.Errorf("create querier: %w", err)
	}
	defer querier.Close()

	names, warnings, err := querier.LabelNames()
	if err != nil {
		return nil, fmt.Errorf("list label names: %w", err)
	}
	_ = warnings

	counts := make([]CardinalityEntry, 0, len(names))
	for _, name := range names {
		vals, _, err := querier.LabelValues(name)
		if err != nil {
			continue
		}
		counts = append(counts, CardinalityEntry{LabelName: name, ValueCount: len(vals)})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].ValueCount > counts[j].ValueCount
	})

	if opts.TopN < len(counts) {
		counts = counts[:opts.TopN]
	}

	// Approximate total series via __name__ cardinality.
	total := 0
	for _, e := range counts {
		if e.LabelName == "__name__" {
			total = e.ValueCount
			break
		}
	}

	return &CardinalityReport{
		TotalSeries: total,
		TopN:        counts,
	}, nil
}
