package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// SliceOptions configures a time-range slice of a snapshot.
type SliceOptions struct {
	Path   string
	Output string
	Start  time.Time
	End    time.Time
}

// SliceResult holds the outcome of a slice operation.
type SliceResult struct {
	SeriesWritten int      `json:"series_written"`
	SamplesKept   int64    `json:"samples_kept"`
	Warnings      []string `json:"warnings,omitempty"`
}

// SliceSnapshot copies only the samples within [Start, End] to Output.
func SliceSnapshot(opts SliceOptions) (*SliceResult, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	if opts.Output == "" {
		return nil, fmt.Errorf("output path is required")
	}
	if !opts.End.IsZero() && opts.Start.After(opts.End) {
		return nil, fmt.Errorf("start must not be after end")
	}

	src, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	dst, err := tsdb.Open(opts.Output, nil, nil, tsdb.DefaultOptions(), nil)
	if err != nil {
		return nil, fmt.Errorf("open destination: %w", err)
	}
	defer dst.Close()

	minT := toMillis(opts.Start)
	maxT := int64(^uint64(0) >> 1)
	if !opts.End.IsZero() {
		maxT = toMillis(opts.End)
	}

	querier, err := src.Querier(minT, maxT)
	if err != nil {
		return nil, fmt.Errorf("create querier: %w", err)
	}
	defer querier.Close()

	ss := querier.Select(false, nil, labels.MustNewMatcher(labels.MatchRegexp, "__name__", ".+"))

	app := dst.Appender(nil)
	result := &SliceResult{}

	for ss.Next() {
		series := ss.At()
		it := series.Iterator(nil)
		var ref uint64
		for it.Next() == 1 {
			t, v := it.At()
			var aerr error
			ref, aerr = app.Append(0, series.Labels(), t, v)
			if aerr != nil {
				result.Warnings = append(result.Warnings, aerr.Error())
				continue
			}
			_ = ref
			result.SamplesKept++
		}
		result.SeriesWritten++
	}
	if err := ss.Err(); err != nil {
		return nil, fmt.Errorf("iterate series: %w", err)
	}
	if err := app.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return result, nil
}
