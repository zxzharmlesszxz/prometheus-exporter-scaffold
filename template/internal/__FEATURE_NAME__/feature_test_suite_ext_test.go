package __FEATURE_NAME__

import (
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestFeatureContract(t *testing.T) {
	suite := NewFeatureTestSuite(NewFeatureTestSpec())
	RegisterFeatureTests(suite)
	suite.RunTests(t)
}

func NewFeatureTestSpec() FeatureTestSpec {
	return FeatureTestSpec{
		SuccessfulSnapshot: func(at time.Time) Snapshot {
			return Snapshot{
				AttemptTime: at,
				Success:     true,
				Value:       1,
			}
		},
		FailedSnapshot: func(at time.Time, err error) Snapshot {
			return Snapshot{
				AttemptTime: at,
				Success:     false,
				Err:         err,
			}
		},
		CheckDefaultSnapshotter: true,
	}
}

func RegisterFeatureTests(suite *FeatureTestSuite) {
	suite.Register("collector_exports_snapshot", testCollectorExportsSnapshot)
	suite.Register("smoke_spec_includes_skeleton_metric", testSmokeSpecIncludesSkeletonMetric)
}

func testCollectorExportsSnapshot(t *testing.T) {
	now := time.Unix(1700000000, 0)
	collector := newTestCollectorWithNow(testFeatureName, testMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), newFakeSnapshotter(Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       42,
	}), testRefreshInterval, func() time.Time { return now })

	expected := fmt.Sprintf(`
# HELP %[1]s Example %[5]s metric emitted by the generated exporter skeleton
# TYPE %[1]s gauge
%[1]s 42
# HELP %[2]s Whether the last %[5]s data collection succeeded
# TYPE %[2]s gauge
%[2]s 1
# HELP %[3]s Unix timestamp of the last %[5]s data collection attempt
# TYPE %[3]s gauge
%[3]s 1.7e+09
# HELP %[4]s Unix timestamp of the last successful %[5]s data collection
# TYPE %[4]s gauge
%[4]s 1.7e+09
`, metricName(testFeatureName, "", metricExampleValue), testLastSuccess, testLastTimestamp, testLastSuccessfulTS, testFeatureName)

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		metricName(testFeatureName, "", metricExampleValue),
		testLastSuccess,
		testLastTimestamp,
		testLastSuccessfulTS,
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func testSmokeSpecIncludesSkeletonMetric(t *testing.T) {
	spec := newTestExporter().SmokeSpec()
	want := metricName(testFeatureName, "", metricExampleValue) + " 1"
	if !hasString(spec.WantMetrics, want) {
		t.Fatalf("SmokeSpec().WantMetrics = %v, want %q", spec.WantMetrics, want)
	}
}
