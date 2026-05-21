package snapshot

import (
	"fmt"
	"sort"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// SeriesInfo holds metadata about a single series in the snapshot.
type SeriesInfo struct {
	Labels     map[string]string `json:"labels"`
	NumSamples int64             `json:"num_samples"`
	MinTime    int64             `json:"min_time_ms"`
	MaxTime    int64             `json:"max_time_ms"`
}

// SeriesSummary holds aggregated series metadata.
type SeriesSummary struct {
	TotalSeries int64        `json:"total_series"`
	TotalSamples int64       `json:"total_samples"`
	Series      []SeriesInfo `json:"series"`
}

// ListSeries returns metadata for all series matching the given matchers.
func ListSeries(db *tsdb.DB, matchers []*labels.Matcher, limit int) (*SeriesSummary, error) {
	if db == nil {
		return nil, fmt.Errorf("db must not be nil")
	}
	if limit <= 0 {
		limit = 100
	}

	q, err := db.Querier(db.Head().MinTime(), db.Head().MaxTime())
	if err != nil {
		return nil, fmt.Errorf("creating querier: %w", err)
	}
	defer q.Close()

	ss := q.Select(false, nil, matchers...)

	var infos []SeriesInfo
	var totalSamples int64

	for ss.Next() {
		if len(infos) >= limit {
			break
		}
		s := ss.At()
		lblMap := make(map[string]string)
		s.Labels().Range(func(l labels.Label) {
			lblMap[l.Name] = l.Value
		})

		var count int64
		var minT, maxT int64 = 1<<62, -1 << 62
		it := s.Iterator(nil)
		for it.Next() != 0 {
			t, _ := it.At()
			if t < minT {
				minT = t
			}
			if t > maxT {
				maxT = t
			}
			count++
		}
		totalSamples += count
		infos = append(infos, SeriesInfo{
			Labels:     lblMap,
			NumSamples: count,
			MinTime:    minT,
			MaxTime:    maxT,
		})
	}
	if err := ss.Err(); err != nil {
		return nil, fmt.Errorf("iterating series: %w", err)
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].NumSamples > infos[j].NumSamples
	})

	return &SeriesSummary{
		TotalSeries:  int64(len(infos)),
		TotalSamples: totalSamples,
		Series:       infos,
	}, nil
}
