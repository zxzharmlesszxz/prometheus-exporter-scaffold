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

func (FeatureExtension) NewSnapshotter(featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return FeatureSnapshotGatherer{}, nil
}

func (FeatureExtension) SmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		WantMetrics: []string{metricExampleValue(ctx.FeatureName) + " 1"},
	}
}
