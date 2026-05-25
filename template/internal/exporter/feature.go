package exporter

import (
	feature "__GO_MODULE__/internal/__FEATURE_NAME__"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

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

func (f *Feature) RegisterFlags(app *kingpin.Application) {
	f.feature.RegisterFlags(app)
}

func (f *Feature) RegisterCollectors(ctx framework.FeatureContext, registry *prometheus.Registry) error {
	return f.feature.RegisterCollectors(ctx, registry)
}

func (f *Feature) RuntimeConfig() []any {
	return f.feature.RuntimeConfig()
}
