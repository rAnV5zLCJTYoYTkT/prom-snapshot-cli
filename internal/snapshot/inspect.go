package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// InspectOptions configures the inspection run.
type InspectOptions struct {
	Path    string
	Verbose bool
}

// BlockInspection holds per-block detail.
type BlockInspection struct {
	ULID       string            `json:"ulid"`
	MinTime    time.Time         `json:"min_time"`
	MaxTime    time.Time         `json:"max_time"`
	Duration   string            `json:"duration"`
	Series     int               `json:"series"`
	Samples    int64             `json:"samples"`
	Chunks     int               `json:"chunks"`
	Labels     map[string]string `json:"labels,omitempty"`
	Compaction int               `json:"compaction_level"`
}

// InspectReport aggregates all block inspections.
type InspectReport struct {
	Path   string            `json:"path"`
	Blocks []BlockInspection `json:"blocks"`
	Total  int               `json:"total_blocks"`
}

// InspectSnapshot opens the TSDB at path and returns detailed block info.
func InspectSnapshot(opts InspectOptions) (*InspectReport, error) {
	if opts.Path == "" {
		return nil, fmt.Errorf("path is required")
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	blocks, err := db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("list blocks: %w", err)
	}

	report := &InspectReport{
		Path:   opts.Path,
		Blocks: make([]BlockInspection, 0, len(blocks)),
		Total:  len(blocks),
	}

	for _, b := range blocks {
		meta := b.Meta()
		minT := time.UnixMilli(meta.MinTime)
		maxT := time.UnixMilli(meta.MaxTime)
		insp := BlockInspection{
			ULID:       meta.ULID.String(),
			MinTime:    minT,
			MaxTime:    maxT,
			Duration:   maxT.Sub(minT).String(),
			Series:     int(meta.Stats.NumSeries),
			Samples:    int64(meta.Stats.NumSamples),
			Chunks:     int(meta.Stats.NumChunks),
			Compaction: meta.Compaction.Level,
		}
		if opts.Verbose {
			insp.Labels = meta.Compaction.Parents != nil && len(meta.Compaction.Parents) > 0 ?
				map[string]string{"parent_count": fmt.Sprintf("%d", len(meta.Compaction.Parents))} : nil
		}
		report.Blocks = append(report.Blocks, insp)
	}

	return report, nil
}
