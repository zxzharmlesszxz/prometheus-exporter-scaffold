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

	collector := NewCollector("", "", nil, newFakeSnapshotter(Snapshot{
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
`, metricLastCollectionSuccess, metricLastCollectionTimestampSeconds, metricLastSuccessfulCollectionTimestampSeconds, defaultFeatureName)

	if err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
		metricExampleValue(defaultFeatureName),
		metricLastCollectionSuccess,
		metricLastCollectionTimestampSeconds,
		metricLastSuccessfulCollectionTimestampSeconds,
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func TestCollectorUsesDefaultSnapshotter(t *testing.T) {
	t.Parallel()

	now := time.Unix(1700000000, 0)
	collector := newCollectorWithNow(defaultFeatureName, defaultMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), nil, time.Minute, func() time.Time {
		return now
	})

	families := exportertest.RegisterAndGather(t, collector)
	exportertest.AssertMetricValue(t, families, metricExampleValue(defaultFeatureName), nil, 1)
	exportertest.AssertMetricValue(t, families, metricLastCollectionSuccess, nil, 1)
	exportertest.AssertMetricValue(t, families, metricLastCollectionTimestampSeconds, nil, float64(now.Unix()))
	exportertest.AssertMetricValue(t, families, metricLastSuccessfulCollectionTimestampSeconds, nil, float64(now.Unix()))
}
