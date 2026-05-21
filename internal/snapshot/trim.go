package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// TrimOptions configures the time range for trimming a snapshot.
type TrimOptions struct {
	Start time.Time
	End   time.Time
}

// TrimResult summarises the outcome of a trim operation.
type TrimResult struct {
	BlocksInspected int     `json:"blocks_inspected"`
	BlocksRetained  int     `json:"blocks_retained"`
	BlocksDropped   int     `json:"blocks_dropped"`
	SeriesRetained  uint64  `json:"series_retained"`
	StartMs         int64   `json:"start_ms"`
	EndMs           int64   `json:"end_ms"`
}

// TrimSnapshot opens the TSDB at snapshotPath and reports which blocks fall
// within [opts.Start, opts.End]. It does NOT delete any data; it only
// describes what would be retained.
func TrimSnapshot(snapshotPath string, opts TrimOptions) (*TrimResult, error) {
	if snapshotPath == "" {
		return nil, fmt.Errorf("snapshot path must not be empty")
	}
	if !opts.End.After(opts.Start) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	db, err := tsdb.OpenDBReadOnly(snapshotPath, "", nil)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	blocks, err := db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("list blocks: %w", err)
	}

	startMs := opts.Start.UnixMilli()
	endMs := opts.End.UnixMilli()

	result := &TrimResult{
		BlocksInspected: len(blocks),
		StartMs:         startMs,
		EndMs:           endMs,
	}

	for _, b := range blocks {
		meta := b.Meta()
		// Retain block if it overlaps with [startMs, endMs]
		if meta.MaxTime > startMs && meta.MinTime < endMs {
			result.BlocksRetained++
			result.SeriesRetained += meta.Stats.NumSeries
		} else {
			result.BlocksDropped++
		}
	}

	return result, nil
}
