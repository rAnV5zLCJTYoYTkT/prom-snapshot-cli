package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

// ExportFormat defines the output format for exported data.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// Sample represents a single time-series data point.
type Sample struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// SeriesExport holds labels and samples for one time series.
type SeriesExport struct {
	Labels  map[string]string `json:"labels"`
	Samples []Sample          `json:"samples"`
}

// ExportOptions configures the export behaviour.
type ExportOptions struct {
	QueryOptions
	Format ExportFormat
	Out    io.Writer
}

// ExportSeries queries the snapshot and writes results to opts.Out.
func ExportSeries(r *Reader, opts ExportOptions) error {
	results, err := QuerySeries(r, opts.QueryOptions)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	var exports []SeriesExport
	for _, ss := range results {
		se := SeriesExport{
			Labels:  labelsToMap(ss.Labels()),
			Samples: []Sample{},
		}
		it := ss.Iterator(nil)
		for it.Next() == chunkenc.ValFloat {
			ts, v := it.At()
			se.Samples = append(se.Samples, Sample{
				Timestamp: time.UnixMilli(ts).UTC(),
				Value:     v,
			})
		}
		if err := it.Err(); err != nil {
			return fmt.Errorf("iterator error: %w", err)
		}
		exports = append(exports, se)
	}

	switch opts.Format {
	case FormatCSV:
		return writeCSV(opts.Out, exports)
	default:
		return writeJSON(opts.Out, exports)
	}
}

func labelsToMap(lbs labels.Labels) map[string]string {
	m := make(map[string]string, lbs.Len())
	lbs.Range(func(l labels.Label) {
		m[l.Name] = l.Value
	})
	return m
}

func writeJSON(w io.Writer, exports []SeriesExport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(exports)
}

func writeCSV(w io.Writer, exports []SeriesExport) error {
	_, err := fmt.Fprintln(w, "labels,timestamp,value")
	if err != nil {
		return err
	}
	for _, se := range exports {
		lblStr := labelsMapToString(se.Labels)
		for _, s := range se.Samples {
			_, err := fmt.Fprintf(w, "%s,%s,%g\n", lblStr, s.Timestamp.Format(time.RFC3339), s.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// labelsMapToString serialises a labels map as a deterministic key=value string
// sorted by label name, e.g. {__name__="up", job="prometheus"}.
func labelsMapToString(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, m[k]))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
