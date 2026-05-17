package snapshot

import (
	"sort"
	"strings"

	"github.com/prometheus/prometheus/model/labels"
)

// LabelMatcher holds a parsed label matcher for filtering series.
type LabelMatcher struct {
	Name  string
	Value string
	Type  string // "=", "!=", "=~", "!~"
}

// ParseMatchers parses a slice of matcher strings like 'job="prometheus"'.
func ParseMatchers(matchers []string) ([]*labels.Matcher, error) {
	parsed := make([]*labels.Matcher, 0, len(matchers))
	for _, m := range matchers {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		matcher, err := labels.ParseMatcher(m)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, matcher)
	}
	return parsed, nil
}

// LabelNames returns a sorted list of unique label names from a set of label sets.
func LabelNames(lsets []labels.Labels) []string {
	seen := make(map[string]struct{})
	for _, lset := range lsets {
		lset.Range(func(l labels.Label) {
			seen[l.Name] = struct{}{}
		})
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// LabelValues returns unique values for a given label name across all label sets.
func LabelValues(lsets []labels.Labels, name string) []string {
	seen := make(map[string]struct{})
	for _, lset := range lsets {
		if v := lset.Get(name); v != "" {
			seen[v] = struct{}{}
		}
	}
	vals := make([]string, 0, len(seen))
	for v := range seen {
		vals = append(vals, v)
	}
	sort.Strings(vals)
	return vals
}
