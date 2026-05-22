package exporter

import (
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
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

	exportertest.WaitForMetricValue(t, registry, metricLastCollectionSuccess, nil, 1)
}

func TestFeatureReportsCollectorRegistrationError(t *testing.T) {
	t.Parallel()

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, NewCollector(defaultMetricNamespace, testFeatureContext().Logger, SnapshotGatherer{}, time.Minute))

	feature := NewFeature()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err == nil {
		t.Fatal("RegisterCollectors() error = nil, want duplicate registration error")
	}
}

func TestFeatureRuntimeConfigNormalizesValues(t *testing.T) {
	t.Parallel()

	feature := &Feature{refreshInterval: -time.Second}
	config := feature.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != defaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, defaultRefreshInterval)
	}
}

func TestFeatureMetadata(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	if feature.FeatureName() != defaultFeatureName {
		t.Fatalf("FeatureName() = %q, want %q", feature.FeatureName(), defaultFeatureName)
	}
	if feature.DefaultListenAddress() != defaultListenAddress {
		t.Fatalf("DefaultListenAddress() = %q, want %q", feature.DefaultListenAddress(), defaultListenAddress)
	}
}
