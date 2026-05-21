package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// BlockMeta holds summary metadata for a single TSDB block.
type BlockMeta struct {
	ULID      string        `json:"ulid"`
	MinTime   time.Time     `json:"min_time"`
	MaxTime   time.Time     `json:"max_time"`
	Duration  time.Duration `json:"duration"`
	NumSeries uint64        `json:"num_series"`
	NumSamples uint64       `json:"num_samples"`
	NumChunks uint64        `json:"num_chunks"`
	SizeBytes int64         `json:"size_bytes"`
}

// CompactSummary holds the result of inspecting all blocks in a snapshot.
type CompactSummary struct {
	Blocks     []BlockMeta `json:"blocks"`
	TotalSeries uint64     `json:"total_series"`
	TotalSamples uint64    `json:"total_samples"`
	TotalSizeBytes int64   `json:"total_size_bytes"`
	BlockCount int         `json:"block_count"`
}

// InspectBlocks opens the TSDB at the given path and returns metadata for
// every block found, along with aggregate totals.
func InspectBlocks(snapshotPath string) (*CompactSummary, error) {
	db, err := tsdb.OpenDBReadOnly(snapshotPath, nil)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	blocks, err := db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("list blocks: %w", err)
	}

	summary := &CompactSummary{
		Blocks: make([]BlockMeta, 0, len(blocks)),
	}

	for _, b := range blocks {
		meta := b.Meta()
		minT := time.UnixMilli(meta.MinTime).UTC()
		maxT := time.UnixMilli(meta.MaxTime).UTC()
		dur := maxT.Sub(minT)

		bm := BlockMeta{
			ULID:       meta.ULID.String(),
			MinTime:    minT,
			MaxTime:    maxT,
			Duration:   dur,
			NumSeries:  meta.Stats.NumSeries,
			NumSamples: meta.Stats.NumSamples,
			NumChunks:  meta.Stats.NumChunks,
			SizeBytes:  b.Size(),
		}

		summary.Blocks = append(summary.Blocks, bm)
		summary.TotalSeries += meta.Stats.NumSeries
		summary.TotalSamples += meta.Stats.NumSamples
		summary.TotalSizeBytes += b.Size()
	}

	summary.BlockCount = len(summary.Blocks)
	return summary, nil
}
