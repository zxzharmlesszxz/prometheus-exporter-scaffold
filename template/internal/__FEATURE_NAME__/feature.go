package __FEATURE_NAME__

import "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"

func NewFeature(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return featurekit.NewSnapshotExtensionFeature(options, featurekit.SnapshotFeatureExtension[Config, Snapshot]{
		DefaultRefreshInterval: DefaultRefreshInterval,
		DefaultConfigFunc:      NewDefaultConfig,
		ConfigFileFunc:         FeatureConfigFile,
		RegisterFlagsFunc:      RegisterFeatureConfigFlags,
		ValidateConfigFunc:     ValidateFeatureConfig,
		ResolveConfigFunc:      ResolveFeatureConfig,
		RuntimeConfigFunc:      FeatureRuntimeConfigEntries,
		NewSnapshotEngineFunc:  NewSnapshotEngine,
		DefaultSnapshotEngine:  NewDefaultSnapshotEngine(),
		MetricsFunc:            NewFeatureMetricSet,
		StatusFunc:             FeatureSnapshotStatus,
		SmokeFunc:              FeatureSmoke,
	})
}
