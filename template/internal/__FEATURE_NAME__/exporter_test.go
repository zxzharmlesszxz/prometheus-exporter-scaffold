package __FEATURE_NAME__

import (
	"testing"
	"time"

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
			"refresh_interval":   30 * time.Second,
			"config_file":        defaultFeatureConfigFile(testFeatureName),
			"config_file_loaded": false,
		},
		DuplicateRegistration:       true,
		LastCollectionSuccessMetric: testLastSuccess,
	})
}

func TestExporterRegistersCollectors(t *testing.T) {
	t.Parallel()

	registry := registerTestFeatureCollectors(t, newTestExporter())
	exportertest.WaitForMetricValue(t, registry, testLastSuccess, nil, 1)
}

func TestContractFeatureDefaults(t *testing.T) {
	t.Parallel()

	exporter := newTestExporterWithOptions(featurekit.SpecOptions{})
	parseExporterFlags(t, exporter, []string{})
	config := exporter.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != DefaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, DefaultRefreshInterval)
	}
	wantConfigFile := defaultFeatureConfigFile("")
	if got := exportertest.RuntimeConfigValue(t, config, "config_file"); got != wantConfigFile {
		t.Fatalf("config_file = %q, want %q", got, wantConfigFile)
	}
	if got := exportertest.RuntimeConfigValue(t, config, "config_file_loaded"); got != false {
		t.Fatalf("config_file_loaded = %v, want false", got)
	}
}

func TestExporterSmokeSpecIncludesSkeletonMetric(t *testing.T) {
	t.Parallel()

	spec := newTestExporter().SmokeSpec()
	want := metricExampleValue(testFeatureName) + " 1"
	if !hasString(spec.WantMetrics, want) {
		t.Fatalf("SmokeSpec().WantMetrics = %v, want %q", spec.WantMetrics, want)
	}
}
