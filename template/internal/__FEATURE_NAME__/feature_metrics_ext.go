package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

func NewFeatureMetricHandlers() featurekit.FeatureMetricHandlers[Snapshot] {
	return featurekit.FeatureMetricHandlers[Snapshot]{
		Collect:  CollectFeatureMetrics,
		LogError: LogFeatureSnapshotError,
	}
}

func CollectFeatureMetrics(ctx featurekit.FeatureMetricsContext[Snapshot], ch chan<- prometheus.Metric, snapshot Snapshot, _ time.Time) {
	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		ctx.Descriptors.Get(metricExampleValue),
		prometheus.GaugeValue,
		snapshot.Value,
	)
}

func LogFeatureSnapshotError(ctx featurekit.FeatureMetricsContext[Snapshot], logger *slog.Logger, snapshot Snapshot) {
	if snapshot.Err != nil {
		logger.Error(ctx.FeatureName+" data collection failed", "err", snapshot.Err)
	}
}
