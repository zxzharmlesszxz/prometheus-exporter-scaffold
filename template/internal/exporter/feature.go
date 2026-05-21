package exporter

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
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

func Main() {
	framework.MainFromProject(NewFeature())
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

func (f *Feature) RegisterCollectors(ctx framework.FeatureContext, registry *prometheus.Registry) error {
	collector := NewCollector(
		ctx.Namespace,
		ctx.Logger,
		SnapshotGatherer{},
		framework.NormalizeDuration(f.refreshInterval, defaultRefreshInterval),
	)
	if err := framework.RegisterAndStartCollectors(context.Background(), registry, collector); err != nil {
		return fmt.Errorf("register __FEATURE_NAME__ collector: %w", err)
	}
	return nil
}

func (f *Feature) RuntimeConfig() []any {
	return []any{
		"refresh_interval", framework.NormalizeDuration(f.refreshInterval, defaultRefreshInterval),
	}
}
