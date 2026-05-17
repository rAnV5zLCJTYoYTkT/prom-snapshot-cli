package snapshot

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb"
)

// Sample represents a single time-series data point.
type Sample struct {
	Timestamp time.Time
	Value     float64
}

// Series represents a labeled time series with its samples.
type Series struct {
	Labels  labels.Labels
	Samples []Sample
}

// QueryOptions holds parameters for querying the snapshot.
type QueryOptions struct {
	Matchers []*labels.Matcher
	Start    time.Time
	End      time.Time
}

// QuerySeries queries the snapshot for time series matching the given options.
func QuerySeries(db *tsdb.DB, opts QueryOptions) ([]Series, error) {
	querier, err := db.Querier(
		context.Background(),
		toMillis(opts.Start),
		toMillis(opts.End),
	)
	if err != nil {
		return nil, fmt.Errorf("creating querier: %w", err)
	}
	defer querier.Close()

	ss := querier.Select(false, nil, opts.Matchers...)

	var results []Series
	for ss.Next() {
		s := ss.At()
		it := s.Iterator(nil)

		var samples []Sample
		for it.Next() != storage.ValNone {
			ts, v := it.At()
			samples = append(samples, Sample{
				Timestamp: time.UnixMilli(ts).UTC(),
				Value:     v,
			})
		}
		if err := it.Err(); err != nil {
			return nil, fmt.Errorf("iterating samples: %w", err)
		}

		results = append(results, Series{
			Labels:  s.Labels(),
			Samples: samples,
		})
	}
	if err := ss.Err(); err != nil {
		return nil, fmt.Errorf("iterating series: %w", err)
	}

	return results, nil
}

func toMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
