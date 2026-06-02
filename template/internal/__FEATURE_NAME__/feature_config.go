package __FEATURE_NAME__

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
	"go.yaml.in/yaml/v2"
)

type FeatureConfigSpec struct {
	enabled           bool
	defaultConfig     Config
	registerFlagsFunc func(*kingpin.Application, featurekit.FlagContext, *Config)
	validateFunc      func(Config) error
	resolveFunc       func(string, Config) (Config, string, bool, error)
	runtimeConfigFunc func(featurekit.RuntimeConfigContext[Config], Config) []any
}

func NewFeatureSpec() FeatureSpec {
	return FeatureSpec{
		refreshInterval: DefaultRefreshInterval,
		config: NewFeatureConfigSpec(
			NewDefaultConfig(),
			RegisterFeatureConfigFlags,
			ResolveFeatureConfig,
			FeatureRuntimeConfigEntries,
		),
		snapshot:    NewFeatureSnapshotSpec(FeatureSnapshotStatus),
		snapshotter: NewFeatureSnapshotterSpec(NewFeatureSnapshotter, NewDefaultSnapshotter()),
		metrics:     NewFeatureMetricsSpec(NewFeatureMetricSet),
		smoke:       NewFeatureSmokeSpec(FeatureSmoke),
	}
}

func NewFeatureConfigSpec(defaultConfig Config, registerFlags func(*kingpin.Application, featurekit.FlagContext, *Config), resolveConfig func(string, Config) (Config, string, bool, error), runtimeConfig func(featurekit.RuntimeConfigContext[Config], Config) []any) FeatureConfigSpec {
	return FeatureConfigSpec{
		enabled:           true,
		defaultConfig:     defaultConfig,
		registerFlagsFunc: registerFlags,
		resolveFunc:       resolveConfig,
		runtimeConfigFunc: runtimeConfig,
	}
}

func (s FeatureConfigSpec) DefaultConfig() Config {
	return s.defaultConfig
}

func (s FeatureConfigSpec) RegisterFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	if !s.enabled {
		return
	}
	app.Flag(
		ctx.FeatureName+".config-file", "YAML config file. If unset, "+defaultFeatureConfigFile(ctx.FeatureName)+" is used when it exists",
	).StringVar(&config.ConfigFile)
	if s.registerFlagsFunc != nil {
		s.registerFlagsFunc(app, ctx, config)
	}
}

func (s FeatureConfigSpec) ValidateConfig(config Config) error {
	if !s.enabled || s.validateFunc == nil {
		return nil
	}
	return s.validateFunc(config)
}

func (s FeatureConfigSpec) ResolveConfig(featureName string, config Config) (Config, string, bool, error) {
	if !s.enabled {
		return config, "", false, nil
	}
	if s.resolveFunc == nil {
		configFile := strings.TrimSpace(config.ConfigFile)
		if configFile == "" {
			configFile = defaultFeatureConfigFile(featureName)
		}
		return config, configFile, false, nil
	}
	return s.resolveFunc(featureName, config)
}

func (s FeatureConfigSpec) RuntimeConfig(ctx featurekit.RuntimeConfigContext[Config]) []any {
	if !s.enabled {
		return nil
	}
	config, configFile, loaded, _ := s.ResolveConfig(ctx.FeatureName, ctx.Config)
	values := []any{
		"config_file", configFile,
		"config_file_loaded", loaded,
	}
	if s.runtimeConfigFunc != nil {
		values = append(values, s.runtimeConfigFunc(ctx, config)...)
	}
	return values
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
