package __FEATURE_NAME__

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
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

func TestFeatureExtensionConfigFileHooks(t *testing.T) {
	t.Parallel()

	config := NewDefaultConfig()
	configFile := writeFeatureConfig(t, "{}\n")
	extension := FeatureExtension{}
	app := kingpin.New("test", "")
	extension.RegisterFlags(app, featurekit.FlagContext{FeatureName: testFeatureName}, &config)
	if _, err := app.Parse([]string{"--" + testFeatureName + ".config-file=" + configFile}); err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}
	if got := config.ConfigFile; got != configFile {
		t.Fatalf("ConfigFile = %q, want %q", got, configFile)
	}
	if got := *FeatureConfigFile(&config); got != configFile {
		t.Fatalf("FeatureConfigFile() = %q, want %q", got, configFile)
	}
	if err := extension.ValidateConfig(config); err != nil {
		t.Fatalf("ValidateConfig() error = %v, want nil", err)
	}
	if _, err := extension.NewSnapshotter(featurekit.CollectorContext[Config]{
		FeatureName:     testFeatureName,
		Framework:       testFeatureContext(),
		Config:          Config{ConfigFile: filepath.Join(t.TempDir(), "missing.yml")},
		RefreshInterval: testRefreshInterval,
	}); err == nil {
		t.Fatal("NewSnapshotter() error = nil, want missing explicit config file error")
	}
}

func TestFeatureConfigFileLoader(t *testing.T) {
	t.Parallel()

	if got := defaultFeatureConfigFile(" custom "); got != filepath.Join("/etc/prometheus", "prometheus-custom-exporter.yml") {
		t.Fatalf("defaultFeatureConfigFile(custom) = %q", got)
	}
	if got := defaultFeatureConfigFile(" "); got != filepath.Join("/etc/prometheus", "prometheus-exporter-exporter.yml") {
		t.Fatalf("defaultFeatureConfigFile(empty) = %q", got)
	}

	missingPath := filepath.Join(t.TempDir(), "missing.yml")
	path, loaded, err := loadFeatureConfigFile(testFeatureName, missingPath, &configFile{})
	if err == nil {
		t.Fatal("loadFeatureConfigFile() error = nil, want missing explicit file error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("loadFeatureConfigFile() error = %v, want os.ErrNotExist", err)
	}
	if path != missingPath || loaded {
		t.Fatalf("loadFeatureConfigFile() path/loaded = %q/%v, want %q/false", path, loaded, missingPath)
	}

	badPath := writeFeatureConfig(t, "unknown: true\n")
	if _, loaded, err := loadFeatureConfigFile(testFeatureName, badPath, &configFile{}); err == nil || loaded {
		t.Fatalf("loadFeatureConfigFile(strict) loaded/error = %v/%v, want false/error", loaded, err)
	}

	configPath := writeFeatureConfig(t, "{}\n")
	path, loaded, err = loadFeatureConfigFile(testFeatureName, " "+configPath+" ", &configFile{})
	if err != nil {
		t.Fatalf("loadFeatureConfigFile(valid) error = %v, want nil", err)
	}
	if path != configPath || !loaded {
		t.Fatalf("loadFeatureConfigFile(valid) path/loaded = %q/%v, want %q/true", path, loaded, configPath)
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

func writeFeatureConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "feature.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return path
}
