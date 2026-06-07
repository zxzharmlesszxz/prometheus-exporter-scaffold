package __FEATURE_NAME__

import "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"

func FeatureSmoke(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		ServerArgs: []string{
			"--" + ctx.FeatureName + ".config-file=../examples/" + DefaultFeatureConfigFileName,
		},
		WantMetrics: []string{featurekit.FeatureMetricName(ctx.FeatureName, featurekit.MetricScopeFeature, metricExampleValue, featureMetricSpecs) + " 1"},
	}
}
