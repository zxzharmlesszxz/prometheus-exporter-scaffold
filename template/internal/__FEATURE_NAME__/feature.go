package __FEATURE_NAME__

import (
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type Feature struct {
	featurekit.FeatureDefaults[Config, Snapshot]
}

func NewFeatureContract() featurekit.FeatureContract[Config, Snapshot] {
	return Feature{}
}

func NewFeature(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return featurekit.NewFeature(featurekit.NewContractSnapshotFeatureSpec[Config, Snapshot](
		options,
		NewFeatureContract(),
	))
}

func (Feature) DefaultRefreshInterval() time.Duration {
	return DefaultRefreshInterval
}

func (Feature) DefaultConfig() Config {
	return NewDefaultConfig()
}

func (Feature) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return NewDefaultSnapshotter()
}
