package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// FilterOptions controls which series and samples are retained.
type FilterOptions struct {
	Path     string
	Matchers []*labels.Matcher
	Start    time.Time
	End      time.Time
	Output   string
}

// FilterResult summarises what was written to the output snapshot.
type FilterResult struct {
	SeriesKept    int    `json:"series_kept"`
	SeriesDropped int    `json:"series_dropped"`
	SamplesKept   int64  `json:"samples_kept"`
	OutputPath    string `json:"output_path"`
}

// FilterSnapshot opens the TSDB at opts.Path, applies label matchers and a
// time window, and writes matching series/samples to a new TSDB at opts.Output.
func FilterSnapshot(opts FilterOptions) (*FilterResult, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	if opts.Output == "" {
		return nil, fmt.Errorf("output path is required")
	}
	if !opts.End.IsZero() && opts.Start.After(opts.End) {
		return nil, fmt.Errorf("start must be before end")
	}

	src, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("open source db: %w", err)
	}
	defer src.Close()

	querier, err := src.Querier(toMillis(opts.Start), toMillis(opts.End))
	if err != nil {
		return nil, fmt.Errorf("create querier: %w", err)
	}
	defer querier.Close()

	matchers := opts.Matchers
	if len(matchers) == 0 {
		matchers = []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "__name__", ".+")}
	}

	dst, err := tsdb.Open(opts.Output, nil, nil, tsdb.DefaultOptions(), nil)
	if err != nil {
		return nil, fmt.Errorf("open destination db: %w", err)
	}
	defer dst.Close()

	app := dst.Appender(nil)

	var kept, dropped int
	var samplesKept int64

	ss := querier.Select(false, nil, matchers...)
	for ss.Next() {
		series := ss.At()
		it := series.Iterator(nil)
		var written int64
		for it.Next() == 1 {
			t, v := it.At()
			_, _ = app.Append(0, series.Labels(), t, v)
			written++
		}
		if written > 0 {
			kept++
			samplesKept += written
		} else {
			dropped++
		}
	}
	if err := ss.Err(); err != nil {
		_ = app.Rollback()
		return nil, fmt.Errorf("iterate series: %w", err)
	}
	if err := app.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &FilterResult{
		SeriesKept:    kept,
		SeriesDropped: dropped,
		SamplesKept:   samplesKept,
		OutputPath:    opts.Output,
	}, nil
}
