package __FEATURE_NAME__

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
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
			DefaultSnapshotter:     NewFeatureContract().DefaultSnapshotter(),
			RefreshInterval:        refreshInterval,
			DefaultRefreshInterval: DefaultRefreshInterval,
			StatusFunc:             NewFeatureContract().SnapshotStatus,
			Now:                    now,
		},
		MetricsFunc: NewFeatureContract().NewMetrics,
	})
}

func TestFeatureExtensionFallbacks(t *testing.T) {
	t.Parallel()

	feature := FeatureExtension{}
	config := Config{}

	if got := feature.DefaultRefreshInterval(); got != 0 {
		t.Fatalf("DefaultRefreshInterval() = %v, want 0", got)
	}
	if got := feature.DefaultConfig(); !reflect.DeepEqual(got, Config{}) {
		t.Fatalf("DefaultConfig() = %#v, want zero config", got)
	}

	feature.RegisterFlags(kingpin.New("test", ""), featurekit.FlagContext{FeatureName: testFeatureName}, &config)
	if err := feature.ValidateConfig(config); err != nil {
		t.Fatalf("ValidateConfig() error = %v", err)
	}
	snapshotter, err := feature.NewSnapshotter(featurekit.CollectorContext[Config]{
		FeatureName: testFeatureName,
		Config:      config,
	})
	if err != nil {
		t.Fatalf("NewSnapshotter() error = %v", err)
	}
	if snapshotter != nil {
		t.Fatalf("NewSnapshotter() = %#v, want nil", snapshotter)
	}
	if got := feature.DefaultSnapshotter(); got != nil {
		t.Fatalf("DefaultSnapshotter() = %#v, want nil", got)
	}
	if got := feature.NewMetrics(featurekit.SnapshotMetricsContext[Snapshot]{
		FeatureName: testFeatureName,
		Namespace:   testMetricNamespace,
	}); got != nil {
		t.Fatalf("NewMetrics() = %#v, want nil", got)
	}

	status := feature.SnapshotStatus(Snapshot{})
	if !status.AttemptTime.IsZero() || status.Success {
		t.Fatalf("SnapshotStatus() = %#v, want zero status", status)
	}
	if got := feature.RuntimeConfig(featurekit.RuntimeConfigContext[Config]{
		FeatureName: testFeatureName,
		Config:      config,
	}); got != nil {
		t.Fatalf("RuntimeConfig() = %#v, want nil", got)
	}
	if got := feature.SmokeSpec(featurekit.SmokeContext[Config]{
		FeatureName: testFeatureName,
		Config:      config,
	}); !reflect.DeepEqual(got, featurekit.SmokeSpec{}) {
		t.Fatalf("SmokeSpec() = %#v, want zero spec", got)
	}
}
