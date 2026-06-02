package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type ExampleMetrics struct {
	featureName      string
	exampleValueDesc *prometheus.Desc
}

func NewFeatureMetricSet(ctx featurekit.SnapshotMetricsContext[Snapshot]) FeatureMetricSet {
	return &ExampleMetrics{
		featureName: ctx.FeatureName,
		exampleValueDesc: prometheus.NewDesc(
			metricExampleValue(ctx.FeatureName),
			"Example "+ctx.FeatureName+" metric emitted by the generated exporter skeleton",
			nil,
			nil,
		),
	}
}

func (m *ExampleMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.exampleValueDesc
}

func (m *ExampleMetrics) Collect(ch chan<- prometheus.Metric, snapshot Snapshot, _ time.Time) {
	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		m.exampleValueDesc,
		prometheus.GaugeValue,
		snapshot.Value,
	)
}

func (m *ExampleMetrics) LogSnapshotError(logger *slog.Logger, snapshot Snapshot) {
	if snapshot.Err != nil {
		logger.Error(m.featureName+" data collection failed", "err", snapshot.Err)
	}
}
