package exporter

import (
	"io"
	"log/slog"
	"testing"
	"time"

	feature "__GO_MODULE__/internal/__FEATURE_NAME__"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestFeatureRegistersAndParsesFlags(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	app := kingpin.New("test", "")
	app.Terminate(func(int) {})
	feature.RegisterFlags(app)

	if _, err := app.Parse([]string{"--" + defaultFeatureName + ".refresh-interval=30s"}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := feature.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != 30*time.Second {
		t.Fatalf("refresh_interval = %v, want %v", got, 30*time.Second)
	}
}

func TestFeatureRegistersCollector(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	registry := prometheus.NewRegistry()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v", err)
	}

	exportertest.WaitForMetricValue(t, registry, metricLastCollectionSuccess, nil, 1)
}

func TestFeatureReportsCollectorRegistrationError(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	feature := NewFeature()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v, want nil", err)
	}

	if err := feature.RegisterCollectors(testFeatureContext(), registry); err == nil {
		t.Fatal("RegisterCollectors() error = nil, want duplicate registration error")
	}
}

func TestFeatureRuntimeConfigNormalizesValues(t *testing.T) {
	t.Parallel()

	exporterFeature := NewFeature()
	config := exporterFeature.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != feature.DefaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, feature.DefaultRefreshInterval)
	}
}

func testFeatureContext() framework.FeatureContext {
	return framework.FeatureContext{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ExporterName: defaultExporterName,
		Namespace:    defaultMetricNamespace,
	}
}
