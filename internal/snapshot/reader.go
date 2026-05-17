package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// BlockInfo holds metadata for a single TSDB block.
type BlockInfo struct {
	ULID      string
	MinTime   time.Time
	MaxTime   time.Time
	NumSeries uint64
}

// SnapshotInfo holds aggregated metadata about a TSDB snapshot.
type SnapshotInfo struct {
	BlockCount  int
	SeriesCount uint64
	MinTime     time.Time
	MaxTime     time.Time
	Blocks      []BlockInfo
}

// Reader wraps a Prometheus TSDB database for read-only access.
type Reader struct {
	db *tsdb.DBReadOnly
}

// NewReader opens a Prometheus TSDB snapshot directory in read-only mode.
func NewReader(path string) (*Reader, error) {
	db, err := tsdb.OpenDBReadOnly(path, nil)
	if err != nil {
		return nil, fmt.Errorf("opening TSDB read-only: %w", err)
	}
	return &Reader{db: db}, nil
}

// Info returns aggregated metadata about the snapshot.
func (r *Reader) Info() (*SnapshotInfo, error) {
	blocks, err := r.db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("listing blocks: %w", err)
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("no blocks found in snapshot")
	}

	info := &SnapshotInfo{
		BlockCount: len(blocks),
		MinTime:    time.UnixMilli(blocks[0].Meta().MinTime).UTC(),
		MaxTime:    time.UnixMilli(blocks[0].Meta().MaxTime).UTC(),
	}

	for _, b := range blocks {
		meta := b.Meta()
		blockMinTime := time.UnixMilli(meta.MinTime).UTC()
		blockMaxTime := time.UnixMilli(meta.MaxTime).UTC()

		if blockMinTime.Before(info.MinTime) {
			info.MinTime = blockMinTime
		}
		if blockMaxTime.After(info.MaxTime) {
			info.MaxTime = blockMaxTime
		}

		info.SeriesCount += meta.Stats.NumSeries
		info.Blocks = append(info.Blocks, BlockInfo{
			ULID:      meta.ULID.String(),
			MinTime:   blockMinTime,
			MaxTime:   blockMaxTime,
			NumSeries: meta.Stats.NumSeries,
		})
	}

	return info, nil
}

// Close releases resources held by the reader.
func (r *Reader) Close() error {
	return r.db.Close()
}
