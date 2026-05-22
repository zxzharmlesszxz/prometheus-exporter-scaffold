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

func (f *Feature) RegisterFlags(app *kingpin.Application) {
	app.Flag(
		defaultFeatureName+".refresh-interval", "How often exporter refreshes "+defaultFeatureName+" data",
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
		return fmt.Errorf("register %s collector: %w", defaultFeatureName, err)
	}
	return nil
}

func (f *Feature) RuntimeConfig() []any {
	return []any{
		"refresh_interval", framework.NormalizeDuration(f.refreshInterval, defaultRefreshInterval),
	}
}
