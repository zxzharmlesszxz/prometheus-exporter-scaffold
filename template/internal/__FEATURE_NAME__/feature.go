package __FEATURE_NAME__

import (
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

// Feature is the standard contract implemented by a concrete exporter feature.
type Feature interface {
	featurekit.FeatureContract[Config, Snapshot]
}

// FeatureExtension carries this exporter's feature-specific method overrides.
type FeatureExtension struct {
	featurekit.FeatureDefaults[Config, Snapshot]
}

var _ Feature = FeatureExtension{}

func NewFeatureContract() Feature {
	return FeatureExtension{}
}

func NewFeature(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return featurekit.NewFeature(featurekit.NewContractSnapshotFeatureSpec[Config, Snapshot](
		options,
		NewFeatureContract(),
	))
}

func (FeatureExtension) DefaultRefreshInterval() time.Duration {
	return DefaultRefreshInterval
}

func (FeatureExtension) DefaultConfig() Config {
	return NewDefaultConfig()
}

func (FeatureExtension) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return NewDefaultSnapshotter()
}
