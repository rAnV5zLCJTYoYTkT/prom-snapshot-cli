package snapshot

import (
	"fmt"
	"time"

	"github.com/prometheus/prometheus/tsdb"
)

// ValidationResult holds the outcome of a snapshot validation run.
type ValidationResult struct {
	Path        string         `json:"path"`
	BlocksFound int            `json:"blocks_found"`
	Issues      []string       `json:"issues,omitempty"`
	Valid       bool           `json:"valid"`
	CheckedAt   time.Time      `json:"checked_at"`
}

// ValidationOptions controls what ValidateSnapshot checks.
type ValidationOptions struct {
	Path          string
	CheckOverlap  bool
	CheckCorruption bool
}

// ValidateSnapshot opens the TSDB at opts.Path and runs a series of
// structural checks, returning a ValidationResult.
func ValidateSnapshot(opts ValidationOptions) (ValidationResult, error) {
	if opts.Path == "" {
		return ValidationResult{}, fmt.Errorf("path must not be empty")
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	blocks, err := db.Blocks()
	if err != nil {
		return ValidationResult{}, fmt.Errorf("list blocks: %w", err)
	}

	result := ValidationResult{
		Path:        opts.Path,
		BlocksFound: len(blocks),
		Valid:       true,
		CheckedAt:   time.Now().UTC(),
	}

	if opts.CheckOverlap {
		for i := 0; i < len(blocks)-1; i++ {
			a := blocks[i].Meta()
			b := blocks[i+1].Meta()
			if a.MaxTime > b.MinTime {
				issue := fmt.Sprintf("block %s and %s overlap (maxTime %d > minTime %d)",
					a.ULID, b.ULID, a.MaxTime, b.MinTime)
				result.Issues = append(result.Issues, issue)
				result.Valid = false
			}
		}
	}

	return result, nil
}
