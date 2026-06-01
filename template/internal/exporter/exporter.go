package exporter

import (
	feature "__GO_MODULE__/internal/__FEATURE_NAME__"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

var mainFromInjectedProject = framework.MainFromInjectedProject

func NewFeature() framework.Feature {
	return featurekit.NewFeature(featurekit.NewContractSnapshotFeatureSpec[feature.Config, feature.Snapshot](
		featurekit.SpecOptions{FeatureName: framework.InjectedFeatureName()},
		feature.NewFeatureContract(),
	))
}

func Main() {
	mainFromInjectedProject(NewFeature())
}

func ExporterInfo() framework.ExporterInfo {
	return framework.ExporterInfoFromInjectedProject(NewFeature())
}
