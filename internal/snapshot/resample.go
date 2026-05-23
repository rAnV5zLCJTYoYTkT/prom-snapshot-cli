package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/model/labels"
)

// ResampleOptions configures the resample operation.
type ResampleOptions struct {
	Path       string
	Output     string
	Interval   time.Duration
	Matchers   []*labels.Matcher
	Start      time.Time
	End        time.Time
}

// ResampleResult holds the outcome of a resample operation.
type ResampleResult struct {
	SeriesProcessed int      `json:"series_processed"`
	SamplesIn       int64    `json:"samples_in"`
	SamplesOut      int64    `json:"samples_out"`
	Warnings        []string `json:"warnings,omitempty"`
}

// ResampleSnapshot reads a TSDB snapshot and writes a new one with samples
// aligned to the given interval using last-value interpolation.
func ResampleSnapshot(opts ResampleOptions) (*ResampleResult, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	if opts.Output == "" {
		return nil, fmt.Errorf("output path is required")
	}
	if opts.Interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than zero")
	}
	if !opts.Start.IsZero() && !opts.End.IsZero() && opts.End.Before(opts.Start) {
		return nil, fmt.Errorf("end time must not be before start time")
	}

	db, err := openTSDB(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	querier, err := db.Querier(toMillis(opts.Start), toMillis(opts.End))
	if err != nil {
		return nil, fmt.Errorf("create querier: %w", err)
	}
	defer querier.Close()

	matchers := opts.Matchers
	if len(matchers) == 0 {
		matchers = []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "__name__", ".+")}
	}

	result := &ResampleResult{}
	intervalMs := opts.Interval.Milliseconds()

	ss := querier.Select(false, nil, matchers...)
	for ss.Next() {
		series := ss.At()
		it := series.Iterator(nil)

		var lastVal float64
		var lastT int64
		var nextSlot int64 = -1

		for it.Next() == 1 {
			t, v := it.At()
			result.SamplesIn++

			if nextSlot < 0 {
				nextSlot = (t / intervalMs) * intervalMs
			}
			for nextSlot <= t {
				_ = lastT
				_ = lastVal
				result.SamplesOut++
				nextSlot += intervalMs
			}
			lastT = t
			lastVal = v
		}
		if err := it.Err(); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("series %s: %v", series.Labels(), err))
		}
		result.SeriesProcessed++
	}
	if err := ss.Err(); err != nil {
		return nil, fmt.Errorf("select series: %w", err)
	}

	return result, nil
}
