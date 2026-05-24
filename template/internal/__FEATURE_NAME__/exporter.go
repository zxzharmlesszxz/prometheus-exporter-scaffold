package __FEATURE_NAME__

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

const (
	defaultFeatureName     = "__FEATURE_NAME__"
	defaultMetricNamespace = "__METRIC_NAMESPACE__"
	DefaultRefreshInterval = time.Minute
)

type Exporter struct {
	featureName            string
	defaultRefreshInterval time.Duration
	refreshInterval        time.Duration
}

type ExporterOptions struct {
	FeatureName            string
	DefaultRefreshInterval time.Duration
}

func NewExporter(options ExporterOptions) *Exporter {
	featureName := options.FeatureName
	if featureName == "" {
		featureName = defaultFeatureName
	}
	defaultRefreshInterval := options.DefaultRefreshInterval
	if defaultRefreshInterval <= 0 {
		defaultRefreshInterval = DefaultRefreshInterval
	}
	return &Exporter{
		featureName:            featureName,
		defaultRefreshInterval: defaultRefreshInterval,
		refreshInterval:        defaultRefreshInterval,
	}
}

func (e *Exporter) RegisterFlags(app *kingpin.Application) {
	app.Flag(
		e.featureName+".refresh-interval", "How often exporter refreshes "+e.featureName+" data",
	).Default(e.defaultRefreshInterval.String()).DurationVar(&e.refreshInterval)
}

func (e *Exporter) RegisterCollectors(ctx framework.FeatureContext, registry *prometheus.Registry) error {
	collector := NewCollector(
		e.featureName,
		ctx.Namespace,
		ctx.Logger,
		SnapshotGatherer{},
		framework.NormalizeDuration(e.refreshInterval, e.defaultRefreshInterval),
	)
	if err := framework.RegisterAndStartCollectors(context.Background(), registry, collector); err != nil {
		return fmt.Errorf("register %s collector: %w", e.featureName, err)
	}
	return nil
}

func (e *Exporter) RuntimeConfig() []any {
	return []any{
		"refresh_interval", framework.NormalizeDuration(e.refreshInterval, e.defaultRefreshInterval),
	}
}
