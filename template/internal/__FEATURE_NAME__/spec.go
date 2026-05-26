package __FEATURE_NAME__

import (
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const DefaultRefreshInterval = time.Minute

type Config struct{}

func NewSpec(options featurekit.SpecOptions) featurekit.FeatureSpec[Config, Snapshot] {
	return featurekit.NewSnapshotFeatureSpec(featurekit.SnapshotFeatureSpec[Config, Snapshot]{
		Options:                options,
		DefaultRefreshInterval: DefaultRefreshInterval,
		Config:                 Config{},
		NewSnapshotterFunc:     newSnapshotter,
		DefaultSnapshotter:     SnapshotGatherer{},
		MetricsFunc:            newMetrics,
		StatusFunc:             snapshotStatus,
		SmokeFunc:              smokeSpec,
	})
}

func newSnapshotter(featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return SnapshotGatherer{}, nil
}

func smokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		WantMetrics: []string{metricExampleValue(ctx.FeatureName) + " 1"},
	}
}
