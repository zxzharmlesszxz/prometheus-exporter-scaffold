package __FEATURE_NAME__

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const (
	testFeatureName      = "__FEATURE_NAME__"
	testMetricNamespace  = "__METRIC_NAMESPACE__"
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
