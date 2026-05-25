package __FEATURE_NAME__

import (
	"time"

	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type Exporter = featurekit.Feature[Config, Snapshot]

type ExporterOptions struct {
	FeatureName            string
	DefaultRefreshInterval time.Duration
}

func NewExporter(options ExporterOptions) *Exporter {
	return featurekit.NewFeature(NewSpec(featurekit.SpecOptions{
		FeatureName:            options.FeatureName,
		DefaultRefreshInterval: options.DefaultRefreshInterval,
	}))
}
