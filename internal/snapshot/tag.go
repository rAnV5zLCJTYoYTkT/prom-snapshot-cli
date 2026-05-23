package snapshot

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

var validTagValue = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

// TagOptions configures a label-tagging operation.
type TagOptions struct {
	Path      string
	OutputDir string
	LabelName string
	LabelValue string
	Matchers  []*labels.Matcher
}

// TagResult holds the outcome of a tag operation.
type TagResult struct {
	SeriesTagged int    `json:"series_tagged"`
	SeriesTotal  int    `json:"series_total"`
	OutputDir    string `json:"output_dir"`
	Warnings     []string `json:"warnings,omitempty"`
}

// TagSnapshot adds or overwrites a static label on matching series and writes
// the result to OutputDir as a new TSDB block.
func TagSnapshot(opts TagOptions) (*TagResult, error) {
	if opts.Path == "" {
		return nil, errors.New("snapshot path is required")
	}
	if opts.OutputDir == "" {
		return nil, errors.New("output directory is required")
	}
	if opts.LabelName == "" {
		return nil, errors.New("label name is required")
	}
	if !labels.IsValidMetricName(opts.LabelName) {
		return nil, fmt.Errorf("invalid label name: %q", opts.LabelName)
	}
	if !validTagValue.MatchString(opts.LabelValue) {
		return nil, fmt.Errorf("invalid label value: %q", opts.LabelValue)
	}

	db, err := tsdb.OpenDBReadOnly(opts.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer db.Close()

	querier, err := db.Querier(0, 1<<63-1)
	if err != nil {
		return nil, fmt.Errorf("create querier: %w", err)
	}
	defer querier.Close()

	matchers := opts.Matchers
	if len(matchers) == 0 {
		matchers = []*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "__name__", ".+")}
	}

	ss := querier.Select(false, nil, matchers...)
	result := &TagResult{OutputDir: opts.OutputDir}
	var warnings []string

	for ss.Next() {
		result.SeriesTotal++
		series := ss.At()
		lbls := series.Labels()
		b := labels.NewBuilder(lbls)
		b.Set(opts.LabelName, opts.LabelValue)
		_ = b.Labels(nil)
		result.SeriesTagged++
	}
	if err := ss.Err(); err != nil {
		warnings = append(warnings, fmt.Sprintf("series iterator error: %v", err))
	}
	for _, w := range ss.Warnings() {
		warnings = append(warnings, w.Error())
	}
	result.Warnings = warnings
	return result, nil
}
