package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/prom-snapshot-cli/internal/snapshot"
)

var (
	replayAddr  string
	replayStart string
	replayEnd   string
)

var replayCmd = &cobra.Command{
	Use:   "replay <snapshot-dir>",
	Short: "Serve a Prometheus-compatible HTTP API backed by a TSDB snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runReplay,
}

func init() {
	replayCmd.Flags().StringVar(&replayAddr, "addr", ":9091", "address to listen on")
	replayCmd.Flags().StringVar(&replayStart, "start", "", "restrict data start (RFC3339 or unix)")
	replayCmd.Flags().StringVar(&replayEnd, "end", "", "restrict data end (RFC3339 or unix)")
	replayCmd.Flags().StringArrayVar(&matcherFlags, "match", nil, "label matchers (repeatable)")
	rootCmd.AddCommand(replayCmd)
}

func runReplay(cmd *cobra.Command, args []string) error {
	matchers, err := snapshot.ParseMatchers(matcherFlags)
	if err != nil {
		return fmt.Errorf("parse matchers: %w", err)
	}

	var start, end time.Time
	if replayStart != "" {
		if start, err = time.Parse(time.RFC3339, replayStart); err != nil {
			return fmt.Errorf("parse --start: %w", err)
		}
	}
	if replayEnd != "" {
		if end, err = time.Parse(time.RFC3339, replayEnd); err != nil {
			return fmt.Errorf("parse --end: %w", err)
		}
	}

	opts := snapshot.ReplayOptions{
		SnapshotDir: args[0],
		Addr:        replayAddr,
		Matchers:    matchers,
		Start:       start,
		End:         end,
	}

	srv, err := snapshot.NewReplayServer(opts)
	if err != nil {
		return err
	}
	defer srv.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintf(os.Stderr, "replay server listening on %s\n", replayAddr)
	return srv.ListenAndServe(ctx)
}
