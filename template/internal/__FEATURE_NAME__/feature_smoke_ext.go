package __FEATURE_NAME__

import "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"

func FeatureSmoke(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		ServerArgs: []string{
			"--" + ctx.FeatureName + ".config-file=../examples/__PROJECT_NAME__.yml",
		},
		WantMetrics: []string{metricName(ctx.FeatureName, "", metricExampleValue) + " 1"},
	}
}
