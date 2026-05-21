package exporter

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

type Snapshot struct {
	AttemptTime time.Time
	Success     bool
	Value       float64
	Err         error
}

type SnapshotGatherer struct{}

func (SnapshotGatherer) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}

type Collector struct {
	*framework.SnapshotCollector[Snapshot]

	exampleValueDesc *prometheus.Desc
}

func NewCollector(namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) *Collector {
	return newCollectorWithNow(namespace, logger, snapshotter, refreshInterval, nil)
}

func newCollectorWithNow(namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration, now func() time.Time) *Collector {
	if namespace == "" {
		namespace = "__METRIC_NAMESPACE__"
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
			"__FEATURE_NAME___example_value",
			"Example __FEATURE_NAME__ metric emitted by the generated exporter skeleton",
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

		LastCollectionSuccessHelp:    "Whether the last __FEATURE_NAME__ data collection succeeded",
		LastCollectionTimestampHelp:  "Unix timestamp of the last __FEATURE_NAME__ data collection attempt",
		LastSuccessfulCollectionHelp: "Unix timestamp of the last successful __FEATURE_NAME__ data collection",
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

func snapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}

func logSnapshotError(logger *slog.Logger, snapshot Snapshot) {
	if snapshot.Err != nil {
		logger.Error("__FEATURE_NAME__ data collection failed", "err", snapshot.Err)
	}
}
