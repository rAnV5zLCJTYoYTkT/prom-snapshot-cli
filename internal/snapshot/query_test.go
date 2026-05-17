package snapshot

import (
	"testing"
	"time"

	"github.com/prometheus/prometheus/model/labels"
)

func TestToMillis(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected int64
	}{
		{
			name:     "unix epoch",
			input:    time.Unix(0, 0).UTC(),
			expected: 0,
		},
		{
			name:     "one second",
			input:    time.Unix(1, 0).UTC(),
			expected: 1000,
		},
		{
			name:     "one millisecond",
			input:    time.UnixMilli(500).UTC(),
			expected: 500,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := toMillis(tc.input)
			if got != tc.expected {
				t.Errorf("toMillis(%v) = %d; want %d", tc.input, got, tc.expected)
			}
		})
	}
}

func TestQueryOptionsMatchers(t *testing.T) {
	matchers, err := ParseMatchers(`{job="prometheus",env=~"prod.*"}`)
	if err != nil {
		t.Fatalf("ParseMatchers error: %v", err)
	}

	opts := QueryOptions{
		Matchers: matchers,
		Start:    time.Now().Add(-1 * time.Hour),
		End:      time.Now(),
	}

	if len(opts.Matchers) != 2 {
		t.Errorf("expected 2 matchers, got %d", len(opts.Matchers))
	}

	found := false
	for _, m := range opts.Matchers {
		if m.Name == "job" && m.Type == labels.MatchEqual && m.Value == "prometheus" {
			found = true
		}
	}
	if !found {
		t.Error("expected matcher job=prometheus not found")
	}

	if opts.End.Before(opts.Start) {
		t.Error("End time should not be before Start time")
	}
}
