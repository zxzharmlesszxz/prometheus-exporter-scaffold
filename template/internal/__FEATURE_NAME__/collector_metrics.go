package __FEATURE_NAME__

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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
