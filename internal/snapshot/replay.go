package snapshot

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
)

// ReplayOptions configures the replay HTTP server.
type ReplayOptions struct {
	SnapshotDir string
	Addr        string
	Matchers    []*labels.Matcher
	Start       time.Time
	End         time.Time
}

// ReplayServer exposes a snapshot as a Prometheus-compatible /api/v1/query_range endpoint.
type ReplayServer struct {
	db   *tsdb.DBReadOnly
	opts ReplayOptions
}

// NewReplayServer opens the snapshot and returns a ReplayServer.
func NewReplayServer(opts ReplayOptions) (*ReplayServer, error) {
	db, err := tsdb.OpenDBReadOnly(opts.SnapshotDir, nil)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	return &ReplayServer{db: db, opts: opts}, nil
}

// Close releases the underlying TSDB handle.
func (s *ReplayServer) Close() error {
	return s.db.Close()
}

// Handler returns an http.Handler that serves the replay API.
func (s *ReplayServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/query_range", s.handleQueryRange)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

// ListenAndServe starts the HTTP server and blocks until ctx is cancelled.
func (s *ReplayServer) ListenAndServe(ctx context.Context) error {
	srv := &http.Server{
		Addr:    s.opts.Addr,
		Handler: s.Handler(),
	}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	select {
	case <-ctx.Done():
		return srv.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}
