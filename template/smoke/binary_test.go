package smoke

import (
	"testing"

	"__GO_MODULE__/internal/exporter/variables"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest/smoketest"
)

const (
	projectName            = variables.DefaultExporterName
	featureName            = variables.DefaultFeatureName
	metricNamespace        = variables.DefaultMetricNamespace
	metricBuildInfo        = metricNamespace + "_build_info"
	metricExampleValue     = featureName + "_example_value"
	metricCollectionStatus = metricNamespace + "_last_collection_success"
)

func TestBinarySmoke(t *testing.T) {
	smoketest.RunBinary(t, smoketest.Config{
		ProjectName:         projectName,
		BuildInfoMetric:     metricBuildInfo,
		ForbiddenUsageNames: []string{metricNamespace},
		RenamedExecutable:   "renamed-" + featureName + "-exporter",
		ServerArgs: func(_ *testing.T, _ string) []string {
			return []string{
				"--" + featureName + ".refresh-interval=100ms",
			}
		},
		WantMetrics: []string{
			metricCollectionStatus + " 1",
			metricExampleValue + " 1",
		},
		RejectMetrics: []string{
			metricCollectionStatus + " 0",
		},
	})
}
