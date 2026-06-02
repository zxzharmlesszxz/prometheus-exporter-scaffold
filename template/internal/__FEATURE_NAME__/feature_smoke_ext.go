package __FEATURE_NAME__

import "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"

func FeatureSmoke(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		WantMetrics: []string{metricExampleValue(ctx.FeatureName) + " 1"},
	}
}
