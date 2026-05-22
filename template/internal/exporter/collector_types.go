package exporter

import (
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

type Collector struct {
	*framework.SnapshotCollector[Snapshot]

	exampleValueDesc *prometheus.Desc
}
