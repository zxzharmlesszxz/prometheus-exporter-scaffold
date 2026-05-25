package __FEATURE_NAME__

import (
	"time"

	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const (
	defaultFeatureName     = "__FEATURE_NAME__"
	defaultMetricNamespace = "__METRIC_NAMESPACE__"
	DefaultRefreshInterval = time.Minute
)

type Exporter = featurekit.Feature[Config, Snapshot]

type ExporterOptions struct {
	FeatureName            string
	DefaultRefreshInterval time.Duration
}

func NewExporter(options ExporterOptions) *Exporter {
	return featurekit.NewFeature(NewSpec(featurekit.SpecOptions{
		FeatureName:             options.FeatureName,
		DefaultFeatureName:      defaultFeatureName,
		DefaultRefreshInterval:  options.DefaultRefreshInterval,
		FallbackRefreshInterval: DefaultRefreshInterval,
	}))
}
