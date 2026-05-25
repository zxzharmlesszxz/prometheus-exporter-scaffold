package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type Collector struct {
	*framework.SnapshotCollector[Snapshot]

	featureName string
	metrics     Metrics
}

func NewCollector(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) *Collector {
	return newCollectorWithNow(featureName, namespace, logger, snapshotter, refreshInterval, nil)
}

func newCollectorWithNow(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration, now func() time.Time) *Collector {
	options := featurekit.ResolveSnapshotCollectorOptions(featurekit.SnapshotCollectorOptions[Snapshot]{
		FeatureName:            featureName,
		Namespace:              namespace,
		Logger:                 logger,
		Snapshotter:            snapshotter,
		DefaultSnapshotter:     SnapshotGatherer{},
		RefreshInterval:        refreshInterval,
		DefaultRefreshInterval: DefaultRefreshInterval,
		StatusFunc:             snapshotStatus,
		Now:                    now,
	})
	collector := &Collector{
		featureName: options.FeatureName,
		metrics:     newMetrics(options.FeatureName, options.Namespace, options.Snapshotter),
	}
	options.DescribeFunc = collector.describeSnapshotMetrics
	options.CollectFunc = collector.collectSnapshotMetrics
	options.ErrorLogFunc = collector.logSnapshotError
	collector.SnapshotCollector = featurekit.NewSnapshotCollector(options)
	return collector
}
