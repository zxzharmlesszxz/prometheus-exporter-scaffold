package __FEATURE_NAME__

import (
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type Config struct{}

const DefaultRefreshInterval = time.Minute

func NewDefaultConfig() Config {
	return Config{}
}

func NewFeatureSpec() FeatureSpec {
	return FeatureSpec{
		RefreshInterval:    DefaultRefreshInterval,
		Config:             NewDefaultConfig(),
		NewSnapshotterFunc: NewFeatureSnapshotter,
		DefaultSnapshotter: NewDefaultSnapshotter(),
		MetricsFunc:        NewFeatureMetrics,
		StatusFunc:         FeatureSnapshotStatus,
		SmokeFunc:          FeatureSmokeSpec,
	}
}

func NewFeatureSnapshotter(featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return FeatureSnapshotGatherer{}, nil
}

func FeatureSmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		WantMetrics: []string{metricExampleValue(ctx.FeatureName) + " 1"},
	}
}
