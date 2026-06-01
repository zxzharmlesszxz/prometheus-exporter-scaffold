package __FEATURE_NAME__

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type FeatureMetrics struct {
	featureName      string
	exampleValueDesc *prometheus.Desc
}

func (Feature) NewMetrics(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	return &FeatureMetrics{
		featureName: ctx.FeatureName,
		exampleValueDesc: prometheus.NewDesc(
			metricExampleValue(ctx.FeatureName),
			"Example "+ctx.FeatureName+" metric emitted by the generated exporter skeleton",
			nil,
			nil,
		),
	}
}

func (m *FeatureMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.exampleValueDesc
}

func (m *FeatureMetrics) Collect(ch chan<- prometheus.Metric, snapshot Snapshot, _ time.Time) {
	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		m.exampleValueDesc,
		prometheus.GaugeValue,
		snapshot.Value,
	)
}
