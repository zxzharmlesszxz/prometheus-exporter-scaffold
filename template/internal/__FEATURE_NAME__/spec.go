package __FEATURE_NAME__

import (
	"log/slog"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const DefaultRefreshInterval = time.Minute

type Config struct{}

func NewSpec(options featurekit.SpecOptions) featurekit.FeatureSpec[Config, Snapshot] {
	defaultRefreshInterval := options.DefaultRefreshInterval
	if defaultRefreshInterval <= 0 {
		defaultRefreshInterval = DefaultRefreshInterval
	}
	fallbackRefreshInterval := options.FallbackRefreshInterval
	if fallbackRefreshInterval <= 0 {
		fallbackRefreshInterval = DefaultRefreshInterval
	}

	return featurekit.FeatureSpec[Config, Snapshot]{
		FeatureName:             options.FeatureName,
		DefaultRefreshInterval:  defaultRefreshInterval,
		FallbackRefreshInterval: fallbackRefreshInterval,
		Config:                  Config{},
		NewSnapshotterFunc:      newSnapshotter,
		NewCollectorFunc:        newCollector,
		SmokeFunc:               smokeSpec,
	}
}

func newSnapshotter(featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return SnapshotGatherer{}, nil
}

func newCollector(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) framework.StartableCollector {
	return NewCollector(featureName, namespace, logger, snapshotter, refreshInterval)
}

func smokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		WantMetrics: []string{metricExampleValue(ctx.FeatureName) + " 1"},
	}
}
