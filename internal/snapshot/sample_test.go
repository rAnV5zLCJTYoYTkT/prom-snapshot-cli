package snapshot

import (
	"testing"
	"time"
)

func TestSampleOptions_EmptyPath(t *testing.T) {
	_, err := FetchSamples(SampleOptions{})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSampleOptions_EndBeforeStart(t *testing.T) {
	now := time.Now()
	_, err := FetchSamples(SampleOptions{
		Path:  "/some/path",
		Start: now,
		End:   now.Add(-time.Hour),
	})
	if err == nil {
		t.Fatal("expected error when end is before start")
	}
}

func TestSampleOptions_InvalidMatcher(t *testing.T) {
	now := time.Now()
	_, err := FetchSamples(SampleOptions{
		Path:     "/some/path",
		Start:    now.Add(-time.Hour),
		End:      now,
		Matchers: []string{"!!invalid!!"},
	})
	if err == nil {
		t.Fatal("expected error for invalid matcher")
	}
}

func TestSampleOptions_InvalidPath(t *testing.T) {
	now := time.Now()
	_, err := FetchSamples(SampleOptions{
		Path:  "/nonexistent/snapshot/path",
		Start: now.Add(-time.Hour),
		End:   now,
	})
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestSamplePoint_Fields(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	sp := SamplePoint{Timestamp: now, Value: 3.14}
	if sp.Timestamp != now {
		t.Errorf("expected timestamp %v, got %v", now, sp.Timestamp)
	}
	if sp.Value != 3.14 {
		t.Errorf("expected value 3.14, got %v", sp.Value)
	}
}

func TestSeriesSamples_Fields(t *testing.T) {
	ss := SeriesSamples{
		Labels:  map[string]string{"__name__": "up", "job": "test"},
		Samples: []SamplePoint{{Timestamp: time.Now(), Value: 1.0}},
	}
	if ss.Labels["__name__"] != "up" {
		t.Errorf("unexpected label value: %v", ss.Labels["__name__"])
	}
	if len(ss.Samples) != 1 {
		t.Errorf("expected 1 sample, got %d", len(ss.Samples))
	}
}
