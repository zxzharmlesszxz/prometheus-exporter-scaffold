package __FEATURE_NAME__

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestCollectorDefaultsAndFailureMetrics(t *testing.T) {
	t.Parallel()

	collector := newTestCollector("", "", newFakeSnapshotter(Snapshot{
		Success: false,
		Err:     errors.New("refresh failed"),
	}), 0)

	expected := fmt.Sprintf(`
# HELP %[1]s Whether the last %[4]s data collection succeeded
# TYPE %[1]s gauge
%[1]s 0
# HELP %[2]s Unix timestamp of the last %[4]s data collection attempt
# TYPE %[2]s gauge
%[2]s 0
# HELP %[3]s Unix timestamp of the last successful %[4]s data collection
# TYPE %[3]s gauge
%[3]s 0
`, "exporter_last_collection_success", "exporter_last_collection_timestamp_seconds", "exporter_last_successful_collection_timestamp_seconds", "exporter")

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		metricName("exporter", "", metricExampleValue),
		"exporter_last_collection_success",
		"exporter_last_collection_timestamp_seconds",
		"exporter_last_successful_collection_timestamp_seconds",
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func TestCollectorUsesDefaultSnapshotter(t *testing.T) {
	t.Parallel()

	now := time.Unix(1700000000, 0)
	collector := newTestCollectorWithNow(testFeatureName, testMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), nil, testRefreshInterval, func() time.Time {
		return now
	})

	families := exportertest.RegisterAndGather(t, collector)
	exportertest.AssertMetricValue(t, families, metricName(testFeatureName, "", metricExampleValue), nil, 1)
	exportertest.AssertMetricValue(t, families, testLastSuccess, nil, 1)
	exportertest.AssertMetricValue(t, families, testLastTimestamp, nil, float64(now.Unix()))
	exportertest.AssertMetricValue(t, families, testLastSuccessfulTS, nil, float64(now.Unix()))
}
