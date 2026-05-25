package __FEATURE_NAME__

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	exampleValueDesc *prometheus.Desc
}

func newMetrics(featureName string) Metrics {
	return Metrics{
		exampleValueDesc: prometheus.NewDesc(
			metricExampleValue(featureName),
			"Example "+featureName+" metric emitted by the generated exporter skeleton",
			nil,
			nil,
		),
	}
}

func (c *Collector) describeSnapshotMetrics(ch chan<- *prometheus.Desc) {
	ch <- c.metrics.exampleValueDesc
}

func (c *Collector) collectSnapshotMetrics(ch chan<- prometheus.Metric, snapshot Snapshot, _ time.Time) {
	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.metrics.exampleValueDesc,
		prometheus.GaugeValue,
		snapshot.Value,
	)
}
