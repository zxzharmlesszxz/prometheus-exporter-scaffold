package __FEATURE_NAME__

import (
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const DefaultRefreshInterval = time.Minute

type Feature struct {
	featurekit.FeatureDefaults[Config, Snapshot]
}

func NewFeatureContract() featurekit.FeatureContract[Config, Snapshot] {
	return Feature{}
}

func (Feature) DefaultRefreshInterval() time.Duration {
	return DefaultRefreshInterval
}

func (Feature) DefaultConfig() Config {
	return Config{}
}

func (Feature) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return FeatureSnapshotGatherer{}
}
