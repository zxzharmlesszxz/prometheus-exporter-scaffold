package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type FeatureMetrics struct {
	set FeatureMetricSet
}

type FeatureMetricSet interface {
	Describe(chan<- *prometheus.Desc)
	Collect(chan<- prometheus.Metric, Snapshot, time.Time)
	LogSnapshotError(*slog.Logger, Snapshot)
}

type FeatureMetricSetFactory func(featurekit.SnapshotMetricsContext[Snapshot]) FeatureMetricSet

type FeatureMetricsSpec struct {
	factory FeatureMetricSetFactory
}

func NewFeatureMetricsSpec(factory FeatureMetricSetFactory) FeatureMetricsSpec {
	return FeatureMetricsSpec{factory: factory}
}

func (s FeatureMetricsSpec) New(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	if s.factory == nil {
		return nil
	}
	set := s.factory(ctx)
	if set == nil {
		return nil
	}
	return NewFeatureMetrics(set)
}

func NewFeatureMetrics(set FeatureMetricSet) FeatureMetrics {
	return FeatureMetrics{set: set}
}

func (m FeatureMetrics) Describe(ch chan<- *prometheus.Desc) {
	if m.set == nil {
		return
	}
	m.set.Describe(ch)
}

func (m FeatureMetrics) Collect(ch chan<- prometheus.Metric, snapshot Snapshot, now time.Time) {
	if m.set == nil {
		return
	}
	m.set.Collect(ch, snapshot, now)
}

func (m FeatureMetrics) LogSnapshotError(logger *slog.Logger, snapshot Snapshot) {
	if m.set == nil {
		return
	}
	m.set.LogSnapshotError(logger, snapshot)
}
