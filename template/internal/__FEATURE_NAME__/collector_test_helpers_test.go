package __FEATURE_NAME__

import (
	"context"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const (
	testFeatureName      = "__FEATURE_NAME__"
	testMetricNamespace  = "__METRIC_NAMESPACE__"
	testExporterName     = "__PROJECT_NAME__"
	testRefreshInterval  = time.Minute
	testLastSuccess      = testMetricNamespace + "_last_collection_success"
	testLastTimestamp    = testMetricNamespace + "_last_collection_timestamp_seconds"
	testLastSuccessfulTS = testMetricNamespace + "_last_successful_collection_timestamp_seconds"
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

func newTestExporter() *featurekit.Feature[Config, Snapshot] {
	return newTestExporterWithOptions(featurekit.SpecOptions{FeatureName: testFeatureName})
}

func newTestExporterWithOptions(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return NewFeature(options)
}

func testFeatureContext() framework.FeatureContext {
	return framework.FeatureContext{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ExporterName: testExporterName,
		Namespace:    testMetricNamespace,
	}
}

func parseExporterFlags(t *testing.T, exporter *featurekit.Feature[Config, Snapshot], args []string) {
	t.Helper()

	exportertest.ParseFeatureFlags(t, exporter, args)
}

func hasString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func startTestCollector(t *testing.T, collector framework.StartableCollector) *prometheus.Registry {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	collector.Start(ctx)

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, collector)
	return registry
}

func registerTestFeatureCollectors(t *testing.T, feature interface {
	RegisterCollectors(framework.FeatureContext, *prometheus.Registry) error
}) *prometheus.Registry {
	t.Helper()

	registry := prometheus.NewRegistry()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v", err)
	}
	return registry
}

func newTestCollector(featureName string, namespace string, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) framework.StartableCollector {
	return newTestCollectorWithNow(featureName, namespace, nil, snapshotter, refreshInterval, nil)
}

func newTestCollectorWithNow(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration, now func() time.Time) framework.StartableCollector {
	return featurekit.NewSnapshotMetricsCollector(featurekit.SnapshotMetricsCollectorOptions[Snapshot]{
		SnapshotCollectorOptions: featurekit.SnapshotCollectorOptions[Snapshot]{
			FeatureName:            featureName,
			Namespace:              namespace,
			Logger:                 logger,
			Snapshotter:            snapshotter,
			DefaultSnapshotter:     NewDefaultSnapshotEngine(),
			RefreshInterval:        refreshInterval,
			DefaultRefreshInterval: DefaultRefreshInterval,
			StatusFunc:             FeatureSnapshotStatus,
			Now:                    now,
		},
		MetricsFunc: NewFeatureMetricSet,
	})
}
