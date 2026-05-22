package exporter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestCollectorExportsSnapshot(t *testing.T) {
	t.Parallel()

	now := time.Unix(1700000000, 0)
	collector := newCollectorWithNow(defaultMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), newFakeSnapshotter(Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       42,
	}), time.Minute, func() time.Time { return now })

	expected := fmt.Sprintf(`
# HELP %[1]s Example %[5]s metric emitted by the generated exporter skeleton
# TYPE %[1]s gauge
%[1]s 42
# HELP %[2]s Whether the last %[5]s data collection succeeded
# TYPE %[2]s gauge
%[2]s 1
# HELP %[3]s Unix timestamp of the last %[5]s data collection attempt
# TYPE %[3]s gauge
%[3]s 1.7e+09
# HELP %[4]s Unix timestamp of the last successful %[5]s data collection
# TYPE %[4]s gauge
%[4]s 1.7e+09
`, metricExampleValue, metricLastCollectionSuccess, metricLastCollectionTimestampSeconds, metricLastSuccessfulCollectionTimestampSeconds, defaultFeatureName)

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		metricExampleValue,
		metricLastCollectionSuccess,
		metricLastCollectionTimestampSeconds,
		metricLastSuccessfulCollectionTimestampSeconds,
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func TestCollectorBackgroundRefreshUpdatesSnapshotOutsideScrape(t *testing.T) {
	t.Parallel()

	start := time.Unix(1700000000, 0)
	snapshotter := newFakeSnapshotter(Snapshot{AttemptTime: start, Success: true, Value: 1})
	collector := NewCollector(defaultMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), snapshotter, 20*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	collector.Start(ctx)

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, collector)
	exportertest.WaitForMetricValue(t, registry, metricExampleValue, nil, 1)

	snapshotter.set(Snapshot{AttemptTime: start.Add(time.Minute), Success: false, Err: errors.New("refresh failed")})
	exportertest.WaitForMetricValue(t, registry, metricLastCollectionSuccess, nil, 0)
}

func TestCollectorDefaultsAndFailureMetrics(t *testing.T) {
	t.Parallel()

	collector := NewCollector("", nil, newFakeSnapshotter(Snapshot{
		Success: false,
		Err:     errors.New("refresh failed"),
	}), 0)

	expected := fmt.Sprintf(`
# HELP %[1]s Whether the last %[4]s data collection succeeded
# TYPE %[1]s gauge
%[1]s 0
# HELP %[2]s Unix timestamp of the last %[4]s data collection attempt
# TYPE %[2]s gauge
%[2]s 0
# HELP %[3]s Unix timestamp of the last successful %[4]s data collection
# TYPE %[3]s gauge
%[3]s 0
`, metricLastCollectionSuccess, metricLastCollectionTimestampSeconds, metricLastSuccessfulCollectionTimestampSeconds, defaultFeatureName)

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		metricExampleValue,
		metricLastCollectionSuccess,
		metricLastCollectionTimestampSeconds,
		metricLastSuccessfulCollectionTimestampSeconds,
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func TestCollectorUsesDefaultSnapshotter(t *testing.T) {
	t.Parallel()

	now := time.Unix(1700000000, 0)
	collector := newCollectorWithNow(defaultMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), nil, time.Minute, func() time.Time {
		return now
	})

	families := exportertest.RegisterAndGather(t, collector)
	exportertest.AssertMetricValue(t, families, metricExampleValue, nil, 1)
	exportertest.AssertMetricValue(t, families, metricLastCollectionSuccess, nil, 1)
	exportertest.AssertMetricValue(t, families, metricLastCollectionTimestampSeconds, nil, float64(now.Unix()))
	exportertest.AssertMetricValue(t, families, metricLastSuccessfulCollectionTimestampSeconds, nil, float64(now.Unix()))
}
