package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// BlockStats holds summary statistics for a single TSDB block.
type BlockStats struct {
	ULID        string
	MinTime     time.Time
	MaxTime     time.Time
	Duration    time.Duration
	NumSeries   uint64
	NumSamples  uint64
	NumChunks   uint64
	SizeBytes   int64
}

// SnapshotStats aggregates statistics across all blocks in a snapshot.
type SnapshotStats struct {
	Blocks      []BlockStats
	TotalSeries uint64
	TotalSamples uint64
	TotalChunks  uint64
	TotalBytes   int64
	EarliestTime time.Time
	LatestTime   time.Time
}

// CollectStats reads block metadata from the TSDB and returns aggregated stats.
func CollectStats(db *tsdb.DBReadOnly) (*SnapshotStats, error) {
	blocks, err := db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("listing blocks: %w", err)
	}

	stats := &SnapshotStats{}
	for i, b := range blocks {
		meta := b.Meta()
		minT := time.UnixMilli(meta.MinTime)
		maxT := time.UnixMilli(meta.MaxTime)

		bs := BlockStats{
			ULID:       meta.ULID.String(),
			MinTime:    minT,
			MaxTime:    maxT,
			Duration:   maxT.Sub(minT),
			NumSeries:  meta.Stats.NumSeries,
			NumSamples: meta.Stats.NumSamples,
			NumChunks:  meta.Stats.NumChunks,
			SizeBytes:  b.Size(),
		}
		stats.Blocks = append(stats.Blocks, bs)
		stats.TotalSeries += bs.NumSeries
		stats.TotalSamples += bs.NumSamples
		stats.TotalChunks += bs.NumChunks
		stats.TotalBytes += bs.SizeBytes

		if i == 0 || minT.Before(stats.EarliestTime) {
			stats.EarliestTime = minT
		}
		if maxT.After(stats.LatestTime) {
			stats.LatestTime = maxT
		}
	}
	return stats, nil
}
