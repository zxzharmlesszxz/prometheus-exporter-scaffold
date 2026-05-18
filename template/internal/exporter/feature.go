package exporter

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	template "github.com/zxzharmlesszxz/prometheus-template-exporter/exporter"
)

const (
	defaultListenAddress   = ":__DEFAULT_PORT__"
	defaultRefreshInterval = time.Minute
)

type Feature struct {
	refreshInterval time.Duration
}

func NewFeature() *Feature {
	return &Feature{
		refreshInterval: defaultRefreshInterval,
	}
}

func (f *Feature) FeatureName() string {
	return "__FEATURE_NAME__"
}

func (f *Feature) DefaultListenAddress() string {
	return defaultListenAddress
}

func (f *Feature) RegisterFlags(app *kingpin.Application) {
	app.Flag(
		"__FEATURE_NAME__.refresh-interval", "How often exporter refreshes __FEATURE_NAME__ data",
	).Default(defaultRefreshInterval.String()).DurationVar(&f.refreshInterval)
}

func (f *Feature) RegisterCollectors(ctx template.FeatureContext, registry *prometheus.Registry) error {
	collector := NewCollector(
		ctx.Namespace,
		ctx.Logger,
		SnapshotGatherer{},
		normalizeDuration(f.refreshInterval, defaultRefreshInterval),
	)
	if err := template.RegisterCollectors(registry, collector); err != nil {
		return fmt.Errorf("register __FEATURE_NAME__ collector: %w", err)
	}
	collector.Start(context.Background())
	return nil
}

func (f *Feature) RuntimeConfig() []any {
	return []any{
		"refresh_interval", normalizeDuration(f.refreshInterval, defaultRefreshInterval),
	}
}

func normalizeDuration(value time.Duration, fallback time.Duration) time.Duration {
	if value <= 0 {
		return fallback
	}
	return value
}
