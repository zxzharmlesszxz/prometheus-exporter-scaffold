package __FEATURE_NAME__

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestCollectorBackgroundRefreshUpdatesSnapshotOutsideScrape(t *testing.T) {
	t.Parallel()

	start := time.Unix(1700000000, 0)
	snapshotter := newFakeSnapshotter(Snapshot{AttemptTime: start, Success: true, Value: 1})
	collector := newTestCollectorWithNow(testFeatureName, testMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), snapshotter, 20*time.Millisecond, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	collector.Start(ctx)

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, collector)
	exportertest.WaitForMetricValue(t, registry, metricExampleValue(testFeatureName), nil, 1)

	snapshotter.set(Snapshot{AttemptTime: start.Add(time.Minute), Success: false, Err: errors.New("refresh failed")})
	exportertest.WaitForMetricValue(t, registry, testLastSuccess, nil, 0)
}
