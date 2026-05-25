package __FEATURE_NAME__

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestExporterRegistersAndParsesFlags(t *testing.T) {
	t.Parallel()

	exporter := NewExporter(ExporterOptions{FeatureName: testFeatureName})
	app := kingpin.New("test", "")
	app.Terminate(func(int) {})
	exporter.RegisterFlags(app)

	if _, err := app.Parse([]string{"--" + testFeatureName + ".refresh-interval=30s"}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := exporter.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != 30*time.Second {
		t.Fatalf("refresh_interval = %v, want %v", got, 30*time.Second)
	}
}

func TestExporterRegistersCollector(t *testing.T) {
	t.Parallel()

	exporter := NewExporter(ExporterOptions{FeatureName: testFeatureName})
	registry := prometheus.NewRegistry()
	if err := exporter.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v", err)
	}

	exportertest.WaitForMetricValue(t, registry, testLastSuccess, nil, 1)
}

func TestExporterReportsCollectorRegistrationError(t *testing.T) {
	t.Parallel()

	exporter := NewExporter(ExporterOptions{FeatureName: testFeatureName})
	registry := prometheus.NewRegistry()
	if err := exporter.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v, want nil", err)
	}

	if err := exporter.RegisterCollectors(testFeatureContext(), registry); err == nil {
		t.Fatal("RegisterCollectors() error = nil, want duplicate registration error")
	}
}

func TestNewExporterDefaults(t *testing.T) {
	t.Parallel()

	exporter := NewExporter(ExporterOptions{})
	config := exporter.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != DefaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, DefaultRefreshInterval)
	}
}

func TestExporterSmokeSpecIncludesSkeletonMetric(t *testing.T) {
	t.Parallel()

	exporter := NewExporter(ExporterOptions{FeatureName: testFeatureName})
	spec := exporter.SmokeSpec()
	want := metricExampleValue(testFeatureName) + " 1"
	for _, metric := range spec.WantMetrics {
		if metric == want {
			return
		}
	}
	t.Fatalf("SmokeSpec().WantMetrics = %v, want %q", spec.WantMetrics, want)
}

func testFeatureContext() framework.FeatureContext {
	return framework.FeatureContext{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ExporterName: "__PROJECT_NAME__",
		Namespace:    testMetricNamespace,
	}
}
