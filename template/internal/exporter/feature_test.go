package exporter

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	template "github.com/zxzharmlesszxz/prometheus-template-exporter/exporter"
)

func TestFeatureRegistersAndParsesFlags(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	app := kingpin.New("test", "")
	app.Terminate(func(int) {})
	feature.RegisterFlags(app)

	if _, err := app.Parse([]string{"--__FEATURE_NAME__.refresh-interval=30s"}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if feature.refreshInterval != 30*time.Second {
		t.Fatalf("refreshInterval = %v, want %v", feature.refreshInterval, 30*time.Second)
	}
}

func TestFeatureRegistersCollector(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	registry := prometheus.NewRegistry()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v", err)
	}

	waitForMetricValue(t, registry, "__METRIC_NAMESPACE___last_collection_success", 1)
}

func TestFeatureRuntimeConfigNormalizesValues(t *testing.T) {
	t.Parallel()

	feature := &Feature{refreshInterval: -time.Second}
	config := feature.RuntimeConfig()
	if got := runtimeConfigValue(t, config, "refresh_interval"); got != defaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, defaultRefreshInterval)
	}
}

func TestFeatureMetadata(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	if feature.FeatureName() != "__FEATURE_NAME__" {
		t.Fatalf("FeatureName() = %q, want %q", feature.FeatureName(), "__FEATURE_NAME__")
	}
	if feature.DefaultListenAddress() != defaultListenAddress {
		t.Fatalf("DefaultListenAddress() = %q, want %q", feature.DefaultListenAddress(), defaultListenAddress)
	}
}

func testFeatureContext() template.FeatureContext {
	return template.FeatureContext{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ExporterName: "__PROJECT_NAME__",
		Namespace:    "__METRIC_NAMESPACE__",
	}
}

func runtimeConfigValue(t *testing.T, config []any, key string) any {
	t.Helper()

	for i := 0; i+1 < len(config); i += 2 {
		if config[i] == key {
			return config[i+1]
		}
	}
	t.Fatalf("missing runtime config key %q in %#v", key, config)
	return nil
}

func waitForMetricValue(t *testing.T, registry *prometheus.Registry, name string, want float64) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for {
		families, err := registry.Gather()
		if err != nil {
			t.Fatalf("Gather() error = %v", err)
		}
		if metricValue(families, name) == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("metric %s did not become %v", name, want)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func metricValue(families []*dto.MetricFamily, name string) float64 {
	for _, family := range families {
		if family.GetName() != name || len(family.GetMetric()) == 0 {
			continue
		}
		return family.GetMetric()[0].GetGauge().GetValue()
	}
	return -1
}
