package exporter

import (
	feature "__GO_MODULE__/internal/__FEATURE_NAME__"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func NewFeature() framework.Feature {
	return feature.NewExporter(feature.ExporterOptions{
		FeatureName: framework.InjectedFeatureName(),
	})
}

func Main() {
	framework.MainFromInjectedProject(NewFeature())
}

func ExporterInfo() framework.ExporterInfo {
	return framework.ExporterInfoFromInjectedProject(NewFeature())
}
