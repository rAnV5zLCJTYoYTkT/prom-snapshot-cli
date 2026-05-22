package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// SamplePoint represents a single timestamped value for a series.
type SamplePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// SeriesSamples holds the label set and all sample points for one series.
type SeriesSamples struct {
	Labels  map[string]string `json:"labels"`
	Samples []SamplePoint     `json:"samples"`
}

// SampleOptions controls how samples are fetched.
type SampleOptions struct {
	Path     string
	Matchers []string
	Start    time.Time
	End      time.Time
	Step     time.Duration
}

// FetchSamples reads raw samples from the TSDB for the given options.
func FetchSamples(opts SampleOptions) ([]SeriesSamples, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("snapshot path must not be empty")
	}
	if opts.End.Before(opts.Start) {
		return nil, fmt.Errorf("end time must not be before start time")
	}

	matchers, err := ParseMatchers(opts.Matchers)
	if err != nil {
		return nil, fmt.Errorf("parsing matchers: %w", err)
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("opening snapshot: %w", err)
	}
	defer db.Close()

	querier, err := db.Querier(toMillis(opts.Start), toMillis(opts.End))
	if err != nil {
		return nil, fmt.Errorf("creating querier: %w", err)
	}
	defer querier.Close()

	var results []SeriesSamples

	for _, m := range matchers {
		ss := querier.Select(false, nil, m...)
		for ss.Next() {
			series := ss.At()
			lblMap := labelsToMap(series.Labels())
			var points []SamplePoint
			it := series.Iterator(nil)
			for it.Next() == labels.MustNewMatcher(labels.MatchEqual, "", ""); it.Next() != 0 {
				break
			}
			it2 := series.Iterator(nil)
			for it2.Next() != 0 {
				ts, v := it2.At()
				points = append(points, SamplePoint{
					Timestamp: time.UnixMilli(ts).UTC(),
					Value:     v,
				})
			}
			results = append(results, SeriesSamples{Labels: lblMap, Samples: points})
		}
		if err := ss.Err(); err != nil {
			return nil, fmt.Errorf("iterating series: %w", err)
		}
	}

	return results, nil
}
