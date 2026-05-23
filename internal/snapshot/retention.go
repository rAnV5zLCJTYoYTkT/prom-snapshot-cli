package snapshot

import (
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// RetentionOptions configures the retention enforcement operation.
type RetentionOptions struct {
	Path      string
	Output    string
	MaxAge    time.Duration
	DryRun    bool
}

// RetentionResult summarises blocks removed or marked for removal.
type RetentionResult struct {
	BlocksScanned int      `json:"blocks_scanned"`
	BlocksRemoved int      `json:"blocks_removed"`
	RemovedULIDs  []string `json:"removed_ulids,omitempty"`
	DryRun        bool     `json:"dry_run"`
}

// EnforceRetention opens a TSDB snapshot and removes (or lists) blocks
// whose maximum time falls before now-maxAge.
func EnforceRetention(opts RetentionOptions) (*RetentionResult, error) {
	if opts.Path == "" {
		return nil, errors.New("path is required")
	}
	if opts.MaxAge <= 0 {
		return nil, errors.New("max-age must be a positive duration")
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	blocks, err := db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("list blocks: %w", err)
	}

	cutoff := time.Now().Add(-opts.MaxAge).UnixMilli()

	result := &RetentionResult{
		BlocksScanned: len(blocks),
		DryRun:        opts.DryRun,
	}

	for _, b := range blocks {
		meta := b.Meta()
		if meta.MaxTime < cutoff {
			result.BlocksRemoved++
			result.RemovedULIDs = append(result.RemovedULIDs, meta.ULID.String())
		}
	}

	return result, nil
}
