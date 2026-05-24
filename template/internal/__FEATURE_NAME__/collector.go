package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func NewCollector(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) *Collector {
	return newCollectorWithNow(featureName, namespace, logger, snapshotter, refreshInterval, nil)
}

func newCollectorWithNow(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration, now func() time.Time) *Collector {
	if featureName == "" {
		featureName = defaultFeatureName
	}
	if namespace == "" {
		namespace = defaultMetricNamespace
	}
	if logger == nil {
		logger = slog.Default()
	}
	if snapshotter == nil {
		snapshotter = SnapshotGatherer{}
	}
	if refreshInterval <= 0 {
		refreshInterval = DefaultRefreshInterval
	}

	collector := &Collector{
		featureName: featureName,
		exampleValueDesc: prometheus.NewDesc(
			metricExampleValue(featureName),
			"Example "+featureName+" metric emitted by the generated exporter skeleton",
			nil,
			nil,
		),
	}
	collector.SnapshotCollector = framework.NewSnapshotCollector(framework.SnapshotCollectorOptions[Snapshot]{
		Namespace:       namespace,
		Logger:          logger,
		Snapshotter:     snapshotter,
		RefreshInterval: refreshInterval,
		StatusFunc:      snapshotStatus,
		DescribeFunc:    collector.describeSnapshotMetrics,
		CollectFunc:     collector.collectSnapshotMetrics,
		ErrorLogFunc:    collector.logSnapshotError,
		Now:             now,

		LastCollectionSuccessHelp:    "Whether the last " + featureName + " data collection succeeded",
		LastCollectionTimestampHelp:  "Unix timestamp of the last " + featureName + " data collection attempt",
		LastSuccessfulCollectionHelp: "Unix timestamp of the last successful " + featureName + " data collection",
	})
	return collector
}
