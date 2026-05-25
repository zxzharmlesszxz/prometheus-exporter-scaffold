package __FEATURE_NAME__

import "testing"

func TestSmokeSpecIncludesSkeletonMetric(t *testing.T) {
	t.Parallel()

	spec := SmokeSpec()
	want := metricExampleValue(defaultFeatureName) + " 1"
	for _, metric := range spec.WantMetrics {
		if metric == want {
			return
		}
	}
	t.Fatalf("SmokeSpec().WantMetrics = %v, want %q", spec.WantMetrics, want)
}
