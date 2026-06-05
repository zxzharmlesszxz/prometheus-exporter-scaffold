package __FEATURE_NAME__

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

const (
	testFeatureName      = "__FEATURE_NAME__"
	testMetricNamespace  = "__METRIC_NAMESPACE__"
	testExporterName     = "__PROJECT_NAME__"
	testRefreshInterval  = time.Minute
	testLastSuccess      = testMetricNamespace + "_last_collection_success"
	testLastTimestamp    = testMetricNamespace + "_last_collection_timestamp_seconds"
	testLastSuccessfulTS = testMetricNamespace + "_last_successful_collection_timestamp_seconds"
)

type FeatureTestFunc func(t *testing.T)

type FeatureTestSpec struct {
	SuccessfulSnapshot                      func(time.Time) Snapshot
	FailedSnapshot                          func(time.Time, error) Snapshot
	ContractFlagArgs                        []string
	ContractRuntimeConfig                   map[string]any
	SkipContractLastCollectionSuccessMetric bool
	CollectorFlagArgs                       []string
	SkipRegisterCollectorsTest              bool
	DefaultRuntimeConfig                    map[string]any
	CheckDefaultSnapshotter                 bool
}

type FeatureTestSuite struct {
	spec  FeatureTestSpec
	tests []featureTest
}

type featureTest struct {
	name string
	run  FeatureTestFunc
}

func NewFeatureTestSuite(spec FeatureTestSpec) *FeatureTestSuite {
	suite := &FeatureTestSuite{spec: spec}
	suite.Register("exporter_contract", suite.testExporterContract)
	if !spec.SkipRegisterCollectorsTest {
		suite.Register("exporter_registers_collectors", suite.testExporterRegistersCollectors)
	}
	suite.Register("contract_feature_defaults", suite.testContractFeatureDefaults)
	suite.Register("feature_config_file_hook", suite.testFeatureConfigFileHook)
	suite.Register("feature_config_file_loader", suite.testFeatureConfigFileLoader)
	suite.Register("smoke_spec_includes_config_file", suite.testSmokeSpecIncludesConfigFile)
	suite.Register("metric_name_contract", suite.testMetricNameContract)
	suite.Register("collector_defaults_and_failure_metrics", suite.testCollectorDefaultsAndFailureMetrics)
	suite.Register("collector_background_refresh_updates_snapshot_outside_scrape", suite.testCollectorBackgroundRefreshUpdatesSnapshotOutsideScrape)
	if spec.CheckDefaultSnapshotter {
		suite.Register("collector_uses_default_snapshotter", suite.testCollectorUsesDefaultSnapshotter)
	}
	return suite
}

func (s *FeatureTestSuite) Register(name string, run FeatureTestFunc) {
	if name == "" {
		panic("feature test name is required")
	}
	if run == nil {
		panic("feature test function is required")
	}
	s.tests = append(s.tests, featureTest{name: name, run: run})
}

func (s *FeatureTestSuite) RunTests(t *testing.T) {
	t.Helper()
	if len(s.tests) == 0 {
		t.Fatal("feature test suite has no registered tests")
	}
	for _, test := range s.tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			test.run(t)
		})
	}
}

func (s *FeatureTestSuite) successfulSnapshot(t *testing.T, at time.Time) Snapshot {
	t.Helper()
	if s.spec.SuccessfulSnapshot == nil {
		t.Fatal("FeatureTestSpec.SuccessfulSnapshot is required")
	}
	return s.spec.SuccessfulSnapshot(at)
}

func (s *FeatureTestSuite) failedSnapshot(t *testing.T, at time.Time, err error) Snapshot {
	t.Helper()
	if s.spec.FailedSnapshot == nil {
		t.Fatal("FeatureTestSpec.FailedSnapshot is required")
	}
	return s.spec.FailedSnapshot(at, err)
}

func mergeRuntimeConfig(base map[string]any, overrides map[string]any) map[string]any {
	merged := make(map[string]any, len(base)+len(overrides))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range overrides {
		merged[key] = value
	}
	return merged
}

func (s *FeatureTestSuite) testExporterContract(t *testing.T) {
	flagArgs := []string{
		"--" + testFeatureName + ".refresh-interval=30s",
	}
	flagArgs = append(flagArgs, s.spec.ContractFlagArgs...)
	wantRuntimeConfig := mergeRuntimeConfig(map[string]any{
		"refresh_interval":   30 * time.Second,
		"config_file":        featurekit.DefaultFeatureConfigFile(testFeatureName),
		"config_file_loaded": false,
	}, s.spec.ContractRuntimeConfig)
	lastCollectionSuccessMetric := testLastSuccess
	if s.spec.SkipContractLastCollectionSuccessMetric {
		lastCollectionSuccessMetric = ""
	}

	exportertest.RunFeatureContract(t, exportertest.FeatureContractConfig{
		NewFeature: func() exportertest.FeatureContractFeature {
			return newTestExporter()
		},
		FeatureContext:              testFeatureContext(),
		FlagArgs:                    flagArgs,
		WantRuntimeConfig:           wantRuntimeConfig,
		DuplicateRegistration:       true,
		LastCollectionSuccessMetric: lastCollectionSuccessMetric,
	})
}

func (s *FeatureTestSuite) testExporterRegistersCollectors(t *testing.T) {
	exporter := newTestExporter()
	parseExporterFlags(t, exporter, s.spec.CollectorFlagArgs)
	registry := registerTestFeatureCollectors(t, exporter)
	exportertest.WaitForMetricValue(t, registry, testLastSuccess, nil, 1)
}

func (s *FeatureTestSuite) testContractFeatureDefaults(t *testing.T) {
	exporter := newTestExporterWithOptions(featurekit.SpecOptions{})
	parseExporterFlags(t, exporter, []string{})
	config := exporter.RuntimeConfig()
	wantRuntimeConfig := mergeRuntimeConfig(map[string]any{
		"refresh_interval":   DefaultRefreshInterval,
		"config_file":        featurekit.DefaultFeatureConfigFile(""),
		"config_file_loaded": false,
	}, s.spec.DefaultRuntimeConfig)
	exportertest.AssertRuntimeConfigValues(t, config, wantRuntimeConfig)
}

func (s *FeatureTestSuite) testFeatureConfigFileHook(t *testing.T) {
	config := NewDefaultConfig()
	configFile := writeFeatureTestConfig(t, "{}\n")
	*FeatureConfigFile(&config) = configFile
	if got := config.ConfigFile; got != configFile {
		t.Fatalf("ConfigFile = %q, want %q", got, configFile)
	}

	exporter := newTestExporter()
	parseExporterFlags(t, exporter, []string{"--" + testFeatureName + ".config-file=" + configFile})
	if got := exportertest.RuntimeConfigValue(t, exporter.RuntimeConfig(), "config_file"); got != configFile {
		t.Fatalf("config_file = %q, want %q", got, configFile)
	}
	if got := exportertest.RuntimeConfigValue(t, exporter.RuntimeConfig(), "config_file_loaded"); got != true {
		t.Fatalf("config_file_loaded = %v, want true", got)
	}
	missingExporter := newTestExporter()
	parseExporterFlags(t, missingExporter, []string{"--" + testFeatureName + ".config-file=" + filepath.Join(t.TempDir(), "missing.yml")})
	if err := missingExporter.RegisterCollectors(testFeatureContext(), prometheus.NewRegistry()); err == nil {
		t.Fatal("RegisterCollectors() error = nil, want missing explicit config file error")
	}
}

func (s *FeatureTestSuite) testFeatureConfigFileLoader(t *testing.T) {
	if got := featurekit.DefaultFeatureConfigFile(" custom "); got != filepath.Join("/etc/prometheus", "prometheus-custom-exporter.yml") {
		t.Fatalf("DefaultFeatureConfigFile(custom) = %q, want default custom path", got)
	}
	if got := featurekit.DefaultFeatureConfigFile(" "); got != filepath.Join("/etc/prometheus", "prometheus-exporter-exporter.yml") {
		t.Fatalf("DefaultFeatureConfigFile(empty) = %q, want default exporter path", got)
	}

	missingPath := filepath.Join(t.TempDir(), "missing.yml")
	path, loaded, err := featurekit.LoadFeatureConfigFile(testFeatureName, missingPath, &configFile{})
	if err == nil {
		t.Fatal("LoadFeatureConfigFile() error = nil, want missing explicit file error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("LoadFeatureConfigFile() error = %v, want os.ErrNotExist", err)
	}
	if path != missingPath || loaded {
		t.Fatalf("LoadFeatureConfigFile() path/loaded = %q/%v, want %q/false", path, loaded, missingPath)
	}

	badPath := writeFeatureTestConfig(t, "unknown: true\n")
	if _, loaded, err := featurekit.LoadFeatureConfigFile(testFeatureName, badPath, &configFile{}); err == nil || loaded {
		t.Fatalf("LoadFeatureConfigFile(strict) loaded/error = %v/%v, want false/error", loaded, err)
	}

	configPath := writeFeatureTestConfig(t, "{}\n")
	path, loaded, err = featurekit.LoadFeatureConfigFile(testFeatureName, " "+configPath+" ", &configFile{})
	if err != nil {
		t.Fatalf("LoadFeatureConfigFile(valid) error = %v, want nil", err)
	}
	if path != configPath || !loaded {
		t.Fatalf("LoadFeatureConfigFile(valid) path/loaded = %q/%v, want %q/true", path, loaded, configPath)
	}
}

func (s *FeatureTestSuite) testSmokeSpecIncludesConfigFile(t *testing.T) {
	spec := newTestExporter().SmokeSpec()
	wantConfig := "--" + testFeatureName + ".config-file=../examples/" + DefaultFeatureConfigFileName
	if !hasString(spec.ServerArgs, wantConfig) {
		t.Fatalf("SmokeSpec().ServerArgs = %v, want %q", spec.ServerArgs, wantConfig)
	}
}

func (s *FeatureTestSuite) testMetricNameContract(t *testing.T) {
	if len(featureMetricSpecs) == 0 {
		t.Fatal("featureMetricSpecs is empty")
	}
	if got := metricName("feature", "namespace", featureMetricSpecs[0].ID); got != featureMetricSpecs[0].MetricName("feature", "namespace") {
		t.Fatalf("metricName(known) = %q, want descriptor spec name", got)
	}
	if got := metricName("feature", "namespace", "missing_metric"); got != "missing_metric" {
		t.Fatalf("metricName(missing) = %q, want missing_metric", got)
	}
}

func (s *FeatureTestSuite) testCollectorDefaultsAndFailureMetrics(t *testing.T) {
	collector := newTestCollector("", "", newFakeSnapshotter(s.failedSnapshot(t, time.Time{}, errors.New("refresh failed"))), 0)

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
		"exporter_last_collection_success",
		"exporter_last_collection_timestamp_seconds",
		"exporter_last_successful_collection_timestamp_seconds",
	); err != nil {
		t.Fatalf("CollectAndCompare() error = %v", err)
	}
}

func (s *FeatureTestSuite) testCollectorBackgroundRefreshUpdatesSnapshotOutsideScrape(t *testing.T) {
	start := time.Unix(1700000000, 0)
	snapshotter := newFakeSnapshotter(s.successfulSnapshot(t, start))
	collector := newTestCollectorWithNow(testFeatureName, testMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), snapshotter, 20*time.Millisecond, nil)

	registry := startTestCollector(t, collector)
	exportertest.WaitForMetricValue(t, registry, testLastSuccess, nil, 1)

	snapshotter.set(s.failedSnapshot(t, start.Add(time.Minute), errors.New("refresh failed")))
	exportertest.WaitForMetricValue(t, registry, testLastSuccess, nil, 0)
}

func (s *FeatureTestSuite) testCollectorUsesDefaultSnapshotter(t *testing.T) {
	now := time.Unix(1700000000, 0)
	collector := newTestCollectorWithNow(testFeatureName, testMetricNamespace, slog.New(slog.NewTextHandler(io.Discard, nil)), nil, testRefreshInterval, func() time.Time {
		return now
	})

	families := exportertest.RegisterAndGather(t, collector)
	exportertest.AssertMetricValue(t, families, testLastSuccess, nil, 1)
	exportertest.AssertMetricValue(t, families, testLastTimestamp, nil, float64(now.Unix()))
	exportertest.AssertMetricValue(t, families, testLastSuccessfulTS, nil, float64(now.Unix()))
}

type fakeSnapshotter struct {
	snapshot atomic.Value
}

func newFakeSnapshotter(snapshot Snapshot) *fakeSnapshotter {
	s := &fakeSnapshotter{}
	s.snapshot.Store(snapshot)
	return s
}

func (s *fakeSnapshotter) Snapshot(context.Context, time.Time) Snapshot {
	return s.snapshot.Load().(Snapshot)
}

func (s *fakeSnapshotter) set(snapshot Snapshot) {
	s.snapshot.Store(snapshot)
}

func newTestExporter() *featurekit.Feature[Config, Snapshot] {
	return newTestExporterWithOptions(featurekit.SpecOptions{FeatureName: testFeatureName})
}

func newTestExporterWithOptions(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return NewFeature(options)
}

func testFeatureContext() framework.FeatureContext {
	return framework.FeatureContext{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ExporterName: testExporterName,
		Namespace:    testMetricNamespace,
	}
}

func parseExporterFlags(t *testing.T, exporter *featurekit.Feature[Config, Snapshot], args []string) {
	t.Helper()

	exportertest.ParseFeatureFlags(t, exporter, args)
}

func hasString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func metricName(featureName string, namespace string, id string) string {
	return featurekit.FeatureMetricName(featureName, namespace, id, featureMetricSpecs)
}

func startTestCollector(t *testing.T, collector framework.StartableCollector) *prometheus.Registry {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	collector.Start(ctx)

	registry := prometheus.NewRegistry()
	exportertest.Register(t, registry, collector)
	return registry
}

func registerTestFeatureCollectors(t *testing.T, feature interface {
	RegisterCollectors(framework.FeatureContext, *prometheus.Registry) error
}) *prometheus.Registry {
	t.Helper()

	registry := prometheus.NewRegistry()
	if err := feature.RegisterCollectors(testFeatureContext(), registry); err != nil {
		t.Fatalf("RegisterCollectors() error = %v", err)
	}
	return registry
}

func newTestCollector(featureName string, namespace string, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration) framework.StartableCollector {
	return newTestCollectorWithNow(featureName, namespace, nil, snapshotter, refreshInterval, nil)
}

func newTestCollectorWithNow(featureName string, namespace string, logger *slog.Logger, snapshotter framework.Snapshotter[Snapshot], refreshInterval time.Duration, now func() time.Time) framework.StartableCollector {
	return featurekit.NewSnapshotMetricsCollector(featurekit.SnapshotMetricsCollectorOptions[Snapshot]{
		SnapshotCollectorOptions: featurekit.SnapshotCollectorOptions[Snapshot]{
			FeatureName:            featureName,
			Namespace:              namespace,
			Logger:                 logger,
			Snapshotter:            snapshotter,
			DefaultSnapshotter:     NewDefaultSnapshotEngine(),
			RefreshInterval:        refreshInterval,
			DefaultRefreshInterval: DefaultRefreshInterval,
			StatusFunc:             FeatureSnapshotStatus,
			Now:                    now,
		},
		MetricsFunc: func(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
			return featurekit.NewFeatureMetrics(ctx, featureMetricSpecs, NewFeatureMetricHandlers())
		},
	})
}

func writeFeatureTestConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "feature.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return path
}
