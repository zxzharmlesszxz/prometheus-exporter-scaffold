package exporter

import feature "__GO_MODULE__/internal/__FEATURE_NAME__"

type Feature struct {
	feature *feature.Exporter
}

func NewFeature() *Feature {
	return &Feature{
		feature: feature.NewExporter(feature.ExporterOptions{
			FeatureName: defaultFeatureName,
		}),
	}
}
