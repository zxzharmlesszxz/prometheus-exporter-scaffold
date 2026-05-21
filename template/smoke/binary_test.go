package smoke

import (
	"testing"

	"github.com/zxzharmlesszxz/prometheus-template-exporter/exporter/exportertest/smoketest"
)

func TestBinarySmoke(t *testing.T) {
	smoketest.RunBinary(t, smoketest.Config{
		ProjectName:         "__PROJECT_NAME__",
		BuildInfoMetric:     "__METRIC_NAMESPACE___build_info",
		ForbiddenUsageNames: []string{"__METRIC_NAMESPACE__"},
		RenamedExecutable:   "renamed-__FEATURE_NAME__-exporter",
		ServerArgs: func(_ *testing.T, _ string) []string {
			return []string{
				"--__FEATURE_NAME__.refresh-interval=100ms",
			}
		},
		WantMetrics: []string{
			"__METRIC_NAMESPACE___last_collection_success 1",
			"__FEATURE_NAME___example_value 1",
		},
		RejectMetrics: []string{
			"__METRIC_NAMESPACE___last_collection_success 0",
		},
	})
}
