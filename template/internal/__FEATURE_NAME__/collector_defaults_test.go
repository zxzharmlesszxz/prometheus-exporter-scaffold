package __FEATURE_NAME__

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
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
		metricExampleValue("exporter"),
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
	exportertest.AssertMetricValue(t, families, metricExampleValue(testFeatureName), nil, 1)
	exportertest.AssertMetricValue(t, families, testLastSuccess, nil, 1)
	exportertest.AssertMetricValue(t, families, testLastTimestamp, nil, float64(now.Unix()))
	exportertest.AssertMetricValue(t, families, testLastSuccessfulTS, nil, float64(now.Unix()))
}

func TestFeatureExtensionFallbacks(t *testing.T) {
	t.Parallel()

	feature := FeatureExtension{}
	config := Config{}

	if got := feature.DefaultRefreshInterval(); got != 0 {
		t.Fatalf("DefaultRefreshInterval() = %v, want 0", got)
	}
	if got := feature.DefaultConfig(); !reflect.DeepEqual(got, Config{}) {
		t.Fatalf("DefaultConfig() = %#v, want zero config", got)
	}

	feature.RegisterFlags(kingpin.New("test", ""), featurekit.FlagContext{FeatureName: testFeatureName}, &config)
	if err := feature.ValidateConfig(config); err != nil {
		t.Fatalf("ValidateConfig() error = %v", err)
	}
	snapshotter, err := feature.NewSnapshotter(featurekit.CollectorContext[Config]{
		FeatureName: testFeatureName,
		Config:      config,
	})
	if err != nil {
		t.Fatalf("NewSnapshotter() error = %v", err)
	}
	if snapshotter != nil {
		t.Fatalf("NewSnapshotter() = %#v, want nil", snapshotter)
	}
	if got := feature.DefaultSnapshotter(); got != nil {
		t.Fatalf("DefaultSnapshotter() = %#v, want nil", got)
	}
	if got := feature.NewMetrics(featurekit.SnapshotMetricsContext[Snapshot]{
		FeatureName: testFeatureName,
		Namespace:   testMetricNamespace,
	}); got != nil {
		t.Fatalf("NewMetrics() = %#v, want nil", got)
	}

	status := feature.SnapshotStatus(Snapshot{})
	if !status.AttemptTime.IsZero() || status.Success {
		t.Fatalf("SnapshotStatus() = %#v, want zero status", status)
	}
	if got := feature.RuntimeConfig(featurekit.RuntimeConfigContext[Config]{
		FeatureName: testFeatureName,
		Config:      config,
	}); got != nil {
		t.Fatalf("RuntimeConfig() = %#v, want nil", got)
	}
	if got := feature.SmokeSpec(featurekit.SmokeContext[Config]{
		FeatureName: testFeatureName,
		Config:      config,
	}); !reflect.DeepEqual(got, featurekit.SmokeSpec{}) {
		t.Fatalf("SmokeSpec() = %#v, want zero spec", got)
	}
}
