package snapshot

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseUnixOrRFC3339_Unix(t *testing.T) {
	t.Parallel()
	got, err := parseUnixOrRFC3339("1609459200")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Unix(1609459200, 0).UTC()
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseUnixOrRFC3339_RFC3339(t *testing.T) {
	t.Parallel()
	got, err := parseUnixOrRFC3339("2021-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseUnixOrRFC3339_Empty(t *testing.T) {
	t.Parallel()
	got, err := parseUnixOrRFC3339("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.IsZero() {
		t.Errorf("expected zero time, got %v", got)
	}
}

func TestHandleQueryRange_BadStart(t *testing.T) {
	t.Parallel()
	srv := &ReplayServer{opts: ReplayOptions{}}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/query_range?start=notadate", nil)
	w := httptest.NewRecorder()
	srv.handleQueryRange(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	var resp apiResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Status != "error" {
		t.Errorf("expected error status, got %q", resp.Status)
	}
}

func TestHealthz(t *testing.T) {
	t.Parallel()
	srv := &ReplayServer{}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	srv.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
