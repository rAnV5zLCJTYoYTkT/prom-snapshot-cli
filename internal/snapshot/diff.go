package snapshot

import (
	"fmt"
	"sort"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// DiffOptions controls the behaviour of DiffSnapshots.
type DiffOptions struct {
	Matchers  []*labels.Matcher
	Start     time.Time
	End       time.Time
	Step      time.Duration
}

// SeriesDiff describes a series that is present in one snapshot but not the
// other, or whose sample count differs.
type SeriesDiff struct {
	Labels  map[string]string `json:"labels"`
	OnlyIn  string            `json:"only_in,omitempty"`
	CountA  int               `json:"count_a,omitempty"`
	CountB  int               `json:"count_b,omitempty"`
}

// DiffResult holds the complete result of comparing two snapshots.
type DiffResult struct {
	Added   []SeriesDiff `json:"added"`
	Removed []SeriesDiff `json:"removed"`
	Changed []SeriesDiff `json:"changed"`
}

// DiffSnapshots compares two open tsdb.DB instances and returns the
// differences in series / sample counts for the given time range.
func DiffSnapshots(dbA, dbB *tsdb.DB, opts DiffOptions) (*DiffResult, error) {
	queryOpts := QueryOptions{
		Matchers: opts.Matchers,
		Start:    opts.Start,
		End:      opts.End,
	}

	seriesA, err := QuerySeries(dbA, queryOpts)
	if err != nil {
		return nil, fmt.Errorf("query snapshot A: %w", err)
	}
	seriesB, err := QuerySeries(dbB, queryOpts)
	if err != nil {
		return nil, fmt.Errorf("query snapshot B: %w", err)
	}

	indexA := indexByFingerprint(seriesA)
	indexB := indexByFingerprint(seriesB)

	result := &DiffResult{}

	for fp, sa := range indexA {
		if sb, ok := indexB[fp]; !ok {
			result.Removed = append(result.Removed, SeriesDiff{
				Labels: labelsToMap(sa.Labels),
				OnlyIn: "A",
			})
		} else if sa.SampleCount != sb.SampleCount {
			result.Changed = append(result.Changed, SeriesDiff{
				Labels: labelsToMap(sa.Labels),
				CountA: sa.SampleCount,
				CountB: sb.SampleCount,
			})
		}
	}
	for fp, sb := range indexB {
		if _, ok := indexA[fp]; !ok {
			result.Added = append(result.Added, SeriesDiff{
				Labels: labelsToMap(sb.Labels),
				OnlyIn: "B",
			})
		}
	}

	sort.Slice(result.Added, func(i, j int) bool { return fmt.Sprint(result.Added[i].Labels) < fmt.Sprint(result.Added[j].Labels) })
	sort.Slice(result.Removed, func(i, j int) bool { return fmt.Sprint(result.Removed[i].Labels) < fmt.Sprint(result.Removed[j].Labels) })
	sort.Slice(result.Changed, func(i, j int) bool { return fmt.Sprint(result.Changed[i].Labels) < fmt.Sprint(result.Changed[j].Labels) })

	return result, nil
}

func indexByFingerprint(series []SeriesResult) map[uint64]SeriesResult {
	m := make(map[uint64]SeriesResult, len(series))
	for _, s := range series {
		m[s.Labels.Hash()] = s
	}
	return m
}
