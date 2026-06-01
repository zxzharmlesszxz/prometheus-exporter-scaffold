package __FEATURE_NAME__

import "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"

func NewSpec(options featurekit.SpecOptions) featurekit.FeatureSpec[Config, Snapshot] {
	return featurekit.NewContractSnapshotFeatureSpec(options, NewFeatureContract())
}
