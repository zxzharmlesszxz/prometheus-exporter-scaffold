package __FEATURE_NAME__

import (
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type Feature = featurekit.Feature[Config, Snapshot]

func NewFeature(options featurekit.SpecOptions) *Feature {
	return featurekit.NewSnapshotExtensionFeature(options, featurekit.SnapshotFeatureExtension[Config, Snapshot]{
		DefaultRefreshInterval: DefaultRefreshInterval,
		DefaultConfigFunc:      NewDefaultConfig,
		ConfigFileFunc:         FeatureConfigFile,
		RegisterFlagsFunc:      RegisterFeatureConfigFlags,
		ValidateConfigFunc:     ValidateFeatureConfig,
		ResolveConfigFunc:      ResolveFeatureConfig,
		RuntimeConfigFunc:      FeatureRuntimeConfigEntries,
		NewSnapshotterFunc:     NewFeatureSnapshotter,
		DefaultSnapshotter:     NewDefaultSnapshotter(),
		MetricsFunc:            NewFeatureMetricSet,
		StatusFunc:             FeatureSnapshotStatus,
		SmokeFunc:              FeatureSmoke,
	})
}
