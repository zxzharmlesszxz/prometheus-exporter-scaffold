package exporter

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

type fakeSnapshotter struct {
	snapshot atomic.Value
}

func newFakeSnapshotter(snapshot Snapshot) *fakeSnapshotter {
	s := &fakeSnapshotter{}
	s.snapshot.Store(snapshot)
	return s
}

func (s *fakeSnapshotter) Snapshot(context.Context, time.Time) Snapshot {
	return s.snapshot.Load().(Snapshot)
}

func (s *fakeSnapshotter) set(snapshot Snapshot) {
	s.snapshot.Store(snapshot)
}

func TestCollectorExportsSnapshot(t *testing.T) {
	t.Parallel()

	now := time.Unix(1700000000, 0)
	collector := newCollectorWithNow("__METRIC_NAMESPACE__", slog.New(slog.NewTextHandler(io.Discard, nil)), newFakeSnapshotter(Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       42,
	}), time.Minute, func() time.Time { return now })

	expected := `
# HELP __FEATURE_NAME___example_value Example __FEATURE_NAME__ metric emitted by the generated exporter skeleton
# TYPE __FEATURE_NAME___example_value gauge
__FEATURE_NAME___example_value 42
# HELP __METRIC_NAMESPACE___last_collection_success Whether the last __FEATURE_NAME__ data collection succeeded
# TYPE __METRIC_NAMESPACE___last_collection_success gauge
__METRIC_NAMESPACE___last_collection_success 1
# HELP __METRIC_NAMESPACE___last_collection_timestamp_seconds Unix timestamp of the last __FEATURE_NAME__ data collection attempt
# TYPE __METRIC_NAMESPACE___last_collection_timestamp_seconds gauge
__METRIC_NAMESPACE___last_collection_timestamp_seconds 1.7e+09
# HELP __METRIC_NAMESPACE___last_successful_collection_timestamp_seconds Unix timestamp of the last successful __FEATURE_NAME__ data collection
# TYPE __METRIC_NAMESPACE___last_successful_collection_timestamp_seconds gauge
__METRIC_NAMESPACE___last_successful_collection_timestamp_seconds 1.7e+09
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		"__FEATURE_NAME___example_value",
		"__METRIC_NAMESPACE___last_collection_success",
		"__METRIC_NAMESPACE___last_collection_timestamp_seconds",
		"__METRIC_NAMESPACE___last_successful_collection_timestamp_seconds",
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func TestCollectorBackgroundRefreshUpdatesSnapshotOutsideScrape(t *testing.T) {
	t.Parallel()

	start := time.Unix(1700000000, 0)
	snapshotter := newFakeSnapshotter(Snapshot{AttemptTime: start, Success: true, Value: 1})
	collector := NewCollector("__METRIC_NAMESPACE__", slog.New(slog.NewTextHandler(io.Discard, nil)), snapshotter, 20*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	collector.Start(ctx)

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, collector)
	exportertest.WaitForMetricValue(t, registry, "__FEATURE_NAME___example_value", nil, 1)

	snapshotter.set(Snapshot{AttemptTime: start.Add(time.Minute), Success: false, Err: errors.New("refresh failed")})
	exportertest.WaitForMetricValue(t, registry, "__METRIC_NAMESPACE___last_collection_success", nil, 0)
}

func TestCollectorDefaultsAndFailureMetrics(t *testing.T) {
	t.Parallel()

	collector := NewCollector("", nil, newFakeSnapshotter(Snapshot{
		Success: false,
		Err:     errors.New("refresh failed"),
	}), 0)

	expected := `
# HELP __METRIC_NAMESPACE___last_collection_success Whether the last __FEATURE_NAME__ data collection succeeded
# TYPE __METRIC_NAMESPACE___last_collection_success gauge
__METRIC_NAMESPACE___last_collection_success 0
# HELP __METRIC_NAMESPACE___last_collection_timestamp_seconds Unix timestamp of the last __FEATURE_NAME__ data collection attempt
# TYPE __METRIC_NAMESPACE___last_collection_timestamp_seconds gauge
__METRIC_NAMESPACE___last_collection_timestamp_seconds 0
# HELP __METRIC_NAMESPACE___last_successful_collection_timestamp_seconds Unix timestamp of the last successful __FEATURE_NAME__ data collection
# TYPE __METRIC_NAMESPACE___last_successful_collection_timestamp_seconds gauge
__METRIC_NAMESPACE___last_successful_collection_timestamp_seconds 0
`

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		"__FEATURE_NAME___example_value",
		"__METRIC_NAMESPACE___last_collection_success",
		"__METRIC_NAMESPACE___last_collection_timestamp_seconds",
		"__METRIC_NAMESPACE___last_successful_collection_timestamp_seconds",
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}
