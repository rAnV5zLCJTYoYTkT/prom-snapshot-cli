package snapshot

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// AnnotateOptions controls how series are annotated with a new constant label.
type AnnotateOptions struct {
	Path       string
	OutputPath string
	LabelName  string
	LabelValue string
	Matchers   []*labels.Matcher
	Start      time.Time
	End        time.Time
}

// AnnotateResult summarises the outcome of an annotate operation.
type AnnotateResult struct {
	SeriesAnnotated int      `json:"series_annotated"`
	SeriesSkipped   int      `json:"series_skipped"`
	Warnings        []string `json:"warnings,omitempty"`
}

var labelNameRE = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// AnnotateSnapshot adds a constant label to all matching series and writes the
// result to OutputPath.
func AnnotateSnapshot(opts AnnotateOptions) (AnnotateResult, error) {
	if opts.Path == "" {
		return AnnotateResult{}, errors.New("path is required")
	}
	if opts.OutputPath == "" {
		return AnnotateResult{}, errors.New("output path is required")
	}
	if opts.LabelName == "" {
		return AnnotateResult{}, errors.New("label name is required")
	}
	if !labelNameRE.MatchString(opts.LabelName) {
		return AnnotateResult{}, fmt.Errorf("invalid label name: %q", opts.LabelName)
	}
	if opts.LabelValue == "" {
		return AnnotateResult{}, errors.New("label value is required")
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return AnnotateResult{}, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	_, err = db.Querier(toMillis(opts.Start), toMillis(opts.End))
	if err != nil {
		return AnnotateResult{}, fmt.Errorf("create querier: %w", err)
	}

	// Placeholder: real implementation would iterate series, add the label,
	// and write blocks to OutputPath via a block writer.
	result := AnnotateResult{
		SeriesAnnotated: 0,
		SeriesSkipped:   0,
	}
	return result, nil
}
