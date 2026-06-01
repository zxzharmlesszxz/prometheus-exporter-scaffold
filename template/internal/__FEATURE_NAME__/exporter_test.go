package __FEATURE_NAME__

import (
	"io"
	"log/slog"
	"testing"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

func TestExporterContract(t *testing.T) {
	t.Parallel()

	exportertest.RunFeatureContract(t, exportertest.FeatureContractConfig{
		NewFeature: func() exportertest.FeatureContractFeature {
			return newTestExporter()
		},
		FeatureContext: testFeatureContext(),
		FlagArgs: []string{
			"--" + testFeatureName + ".refresh-interval=30s",
		},
		WantRuntimeConfig: map[string]any{
			"refresh_interval": 30 * time.Second,
		},
		DuplicateRegistration:       true,
		LastCollectionSuccessMetric: testLastSuccess,
	})
}

func TestContractFeatureDefaults(t *testing.T) {
	t.Parallel()

	exporter := newTestExporterWithOptions(featurekit.SpecOptions{})
	config := exporter.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != DefaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, DefaultRefreshInterval)
	}
}

func TestExporterSmokeSpecIncludesSkeletonMetric(t *testing.T) {
	t.Parallel()

	spec := newTestExporter().SmokeSpec()
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
