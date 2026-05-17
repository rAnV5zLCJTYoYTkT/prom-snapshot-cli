package snapshot

import (
	"time"
)

// mustParseRFC3339 parses an RFC3339 time string and panics on error.
// Intended for use in tests only.
func mustParseRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("mustParseRFC3339: " + err.Error())
	}
	return t
}

// newTestSeriesExport creates a SeriesExport for use in tests.
func newTestSeriesExport(name string, samples []Sample) SeriesExport {
	return SeriesExport{
		Labels: map[string]string{
			"__name__": name,
			"job":      "test",
		},
		Samples: samples,
	}
}

// newTestSamples generates n evenly-spaced samples starting at start.
func newTestSamples(start time.Time, n int, step time.Duration, value float64) []Sample {
	samples := make([]Sample, n)
	for i := 0; i < n; i++ {
		samples[i] = Sample{
			Timestamp: start.Add(time.Duration(i) * step),
			Value:     value,
		}
	}
	return samples
}
