package exporter

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func NewCollector(namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) *Collector {
	return newCollectorWithNow(namespace, logger, snapshotter, refreshInterval, nil)
}

func newCollectorWithNow(namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration, now func() time.Time) *Collector {
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
		refreshInterval = defaultRefreshInterval
	}

	collector := &Collector{
		exampleValueDesc: prometheus.NewDesc(
			metricExampleValue,
			"Example "+defaultFeatureName+" metric emitted by the generated exporter skeleton",
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
		ErrorLogFunc:    logSnapshotError,
		Now:             now,

		LastCollectionSuccessHelp:    "Whether the last " + defaultFeatureName + " data collection succeeded",
		LastCollectionTimestampHelp:  "Unix timestamp of the last " + defaultFeatureName + " data collection attempt",
		LastSuccessfulCollectionHelp: "Unix timestamp of the last successful " + defaultFeatureName + " data collection",
	})
	return collector
}

func (c *Collector) describeSnapshotMetrics(ch chan<- *prometheus.Desc) {
	ch <- c.exampleValueDesc
}

func (c *Collector) collectSnapshotMetrics(ch chan<- prometheus.Metric, snapshot Snapshot, _ time.Time) {
	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.exampleValueDesc,
		prometheus.GaugeValue,
		snapshot.Value,
	)
}
