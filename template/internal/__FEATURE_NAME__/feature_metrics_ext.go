package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type ExampleMetrics struct {
	featureName string
	metrics     metricDescriptors
}

func NewFeatureMetricSet(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	return &ExampleMetrics{
		featureName: ctx.FeatureName,
		metrics:     loadMetricDescriptors(ctx.FeatureName, ctx.Namespace, featureMetricSpecs),
	}
}

func (m *ExampleMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.metrics.Describe(ch)
}

func (m *ExampleMetrics) Collect(ch chan<- prometheus.Metric, snapshot Snapshot, _ time.Time) {
	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		m.metrics.Get(metricExampleValue),
		prometheus.GaugeValue,
		snapshot.Value,
	)
}

func (m *ExampleMetrics) LogSnapshotError(logger *slog.Logger, snapshot Snapshot) {
	if snapshot.Err != nil {
		logger.Error(m.featureName+" data collection failed", "err", snapshot.Err)
	}
}
