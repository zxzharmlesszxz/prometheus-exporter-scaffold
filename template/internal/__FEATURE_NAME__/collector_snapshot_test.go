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

func TestCollectorExportsSnapshot(t *testing.T) {
	t.Parallel()

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
