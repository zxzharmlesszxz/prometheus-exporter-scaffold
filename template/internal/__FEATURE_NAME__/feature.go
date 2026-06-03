package __FEATURE_NAME__

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
	"go.yaml.in/yaml/v2"
)

type Feature interface {
	featurekit.FeatureContract[Config, Snapshot]
}

type FeatureExtension struct {
	featurekit.FeatureDefaults[Config, Snapshot]
}

var _ Feature = FeatureExtension{}

func NewFeatureContract() Feature {
	return FeatureExtension{}
}

func NewFeature(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return featurekit.NewFeature(featurekit.NewContractSnapshotFeatureSpec[Config, Snapshot](
		options,
		NewFeatureContract(),
	))
}

func (FeatureExtension) DefaultRefreshInterval() time.Duration {
	return DefaultRefreshInterval
}

func (FeatureExtension) DefaultConfig() Config {
	return NewDefaultConfig()
}

func (FeatureExtension) RegisterFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	app.Flag(
		ctx.FeatureName+".config-file", "YAML config file. If unset, "+defaultFeatureConfigFile(ctx.FeatureName)+" is used when it exists",
	).StringVar(&config.ConfigFile)
	RegisterFeatureConfigFlags(app, ctx, config)
}

func (FeatureExtension) ValidateConfig(config Config) error {
	return ValidateFeatureConfig(config)
}

func (FeatureExtension) NewSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return NewFeatureSnapshotter(ctx)
}

func (FeatureExtension) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return NewDefaultSnapshotter()
}

func (FeatureExtension) NewMetrics(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	return NewFeatureMetricSet(ctx)
}

func (FeatureExtension) SnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return FeatureSnapshotStatus(snapshot)
}

func (FeatureExtension) RuntimeConfig(ctx featurekit.RuntimeConfigContext[Config]) []any {
	config, configFile, loaded, _ := ResolveFeatureConfig(ctx.FeatureName, ctx.Config)
	values := []any{
		"config_file", configFile,
		"config_file_loaded", loaded,
	}
	return append(values, FeatureRuntimeConfigEntries(ctx, config)...)
}

func (FeatureExtension) SmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return FeatureSmoke(ctx)
}

func defaultFeatureConfigFile(featureName string) string {
	name := strings.TrimSpace(featureName)
	if name == "" {
		name = "exporter"
	}
	return filepath.Join("/etc/prometheus", "prometheus-"+name+"-exporter.yml")
}

func loadFeatureConfigFile(featureName string, explicitPath string, target any) (string, bool, error) {
	path := strings.TrimSpace(explicitPath)
	required := path != ""
	if path == "" {
		path = defaultFeatureConfigFile(featureName)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if !required && errors.Is(err, os.ErrNotExist) {
			return path, false, nil
		}
		return path, false, fmt.Errorf("read %s config file %q: %w", featureName, path, err)
	}
	if err := yaml.UnmarshalStrict(data, target); err != nil {
		return path, false, fmt.Errorf("parse %s config file %q: %w", featureName, path, err)
	}
	return path, true, nil
}
