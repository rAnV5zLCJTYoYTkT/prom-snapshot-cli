package snapshot

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
)

func TestParseMatchers_Valid(t *testing.T) {
	matchers, err := ParseMatchers([]string{`job="prometheus"`, `instance=~"localhost.*"`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matchers) != 2 {
		t.Fatalf("expected 2 matchers, got %d", len(matchers))
	}
}

func TestParseMatchers_Empty(t *testing.T) {
	matchers, err := ParseMatchers([]string{"  ", ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matchers) != 0 {
		t.Fatalf("expected 0 matchers, got %d", len(matchers))
	}
}

func TestParseMatchers_Invalid(t *testing.T) {
	_, err := ParseMatchers([]string{"not_valid_matcher"})
	if err == nil {
		t.Fatal("expected error for invalid matcher, got nil")
	}
}

func TestParseMatchers_Nil(t *testing.T) {
	matchers, err := ParseMatchers(nil)
	if err != nil {
		t.Fatalf("unexpected error for nil input: %v", err)
	}
	if len(matchers) != 0 {
		t.Fatalf("expected 0 matchers for nil input, got %d", len(matchers))
	}
}

func TestLabelNames(t *testing.T) {
	lsets := []labels.Labels{
		labels.FromStrings("job", "prometheus", "instance", "localhost:9090"),
		labels.FromStrings("job", "node", "env", "prod"),
	}
	names := LabelNames(lsets)
	expected := []string{"env", "instance", "job"}
	if len(names) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, names)
	}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("expected %s at index %d, got %s", expected[i], i, n)
		}
	}
}

func TestLabelValues(t *testing.T) {
	lsets := []labels.Labels{
		labels.FromStrings("job", "prometheus"),
		labels.FromStrings("job", "node"),
		labels.FromStrings("job", "prometheus"),
	}
	vals := LabelValues(lsets, "job")
	if len(vals) != 2 {
		t.Fatalf("expected 2 unique values, got %d: %v", len(vals), vals)
	}
	if vals[0] != "node" || vals[1] != "prometheus" {
		t.Errorf("unexpected values: %v", vals)
	}
}
