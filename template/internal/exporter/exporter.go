package exporter

import (
	feature "__GO_MODULE__/internal/__FEATURE_NAME__"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

var mainFromInjectedProject = framework.MainFromInjectedProject

func NewFeature() framework.Feature {
	return feature.NewExporter(feature.ExporterOptions{
		FeatureName: framework.InjectedFeatureName(),
	})
}

func Main() {
	mainFromInjectedProject(NewFeature())
}

func ExporterInfo() framework.ExporterInfo {
	return framework.ExporterInfoFromInjectedProject(NewFeature())
}
