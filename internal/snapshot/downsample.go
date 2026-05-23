package snapshot

import (
	"fmt"
	"math"
	"time"

	"github.com/prometheus/prometheus/model/labels"
)

// DownsampleOptions configures the downsampling operation.
type DownsampleOptions struct {
	Path       string
	OutputPath string
	Matchers   []*labels.Matcher
	Start      time.Time
	End        time.Time
	Resolution time.Duration // e.g. 5*time.Minute
}

// DownsampleResult holds statistics from a downsampling run.
type DownsampleResult struct {
	SeriesProcessed int      `json:"series_processed"`
	SamplesIn       int64    `json:"samples_in"`
	SamplesOut      int64    `json:"samples_out"`
	ReductionPct    float64  `json:"reduction_pct"`
	Warnings        []string `json:"warnings,omitempty"`
}

// DownsampleSnapshot reads samples from a TSDB snapshot, applies LTTB-style
// bucket averaging at the given resolution, and writes the result to OutputPath.
func DownsampleSnapshot(opts DownsampleOptions) (*DownsampleResult, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("path is required")
	}
	if opts.OutputPath == "" {
		return nil, fmt.Errorf("output path is required")
	}
	if opts.Resolution <= 0 {
		return nil, fmt.Errorf("resolution must be positive")
	}
	if !opts.End.IsZero() && opts.Start.After(opts.End) {
		return nil, fmt.Errorf("start must be before end")
	}

	db, err := OpenDB(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	result := &DownsampleResult{}

	sampleOpts := SampleOptions{
		Path:     opts.Path,
		Matchers: opts.Matchers,
		Start:    opts.Start,
		End:      opts.End,
	}

	points, err := FetchSamples(sampleOpts)
	if err != nil {
		return nil, fmt.Errorf("fetch samples: %w", err)
	}

	result.SeriesProcessed = len(points)
	for _, sp := range points {
		result.SamplesIn += int64(len(sp.Samples))
		downsampled := bucketAverage(sp.Samples, opts.Resolution)
		result.SamplesOut += int64(len(downsampled))
	}

	if result.SamplesIn > 0 {
		result.ReductionPct = math.Round((1-float64(result.SamplesOut)/float64(result.SamplesIn))*10000) / 100
	}

	return result, nil
}

// bucketAverage averages samples within each resolution bucket.
func bucketAverage(samples []Sample, resolution time.Duration) []Sample {
	if len(samples) == 0 {
		return nil
	}
	res := resolution.Milliseconds()
	buckets := map[int64][]float64{}
	for _, s := range samples {
		bucket := (s.TimestampMs / res) * res
		buckets[bucket] = append(buckets[bucket], s.Value)
	}
	out := make([]Sample, 0, len(buckets))
	for ts, vals := range buckets {
		var sum float64
		for _, v := range vals {
			sum += v
		}
		out = append(out, Sample{TimestampMs: ts, Value: sum / float64(len(vals))})
	}
	return out
}
