package __FEATURE_NAME__

import (
	"time"

	"github.com/alecthomas/kingpin/v2"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type Config struct {
	ConfigFile string
}

const DefaultRefreshInterval = time.Minute

type configFile struct{}

func NewDefaultConfig() Config {
	return Config{}
}

func NewFeatureSpec() FeatureSpec {
	return FeatureSpec{
		RefreshInterval:    DefaultRefreshInterval,
		Config:             NewDefaultConfig(),
		RegisterFlagsFunc:  RegisterFeatureFlags,
		NewSnapshotterFunc: NewFeatureSnapshotter,
		DefaultSnapshotter: NewDefaultSnapshotter(),
		MetricsFunc:        NewFeatureMetrics,
		StatusFunc:         FeatureSnapshotStatus,
		SmokeFunc:          FeatureSmokeSpec,
	}
}

func RegisterFeatureFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	app.Flag(
		ctx.FeatureName+".config-file", "YAML config file. If unset, "+defaultFeatureConfigFile(ctx.FeatureName)+" is used when it exists",
	).StringVar(&config.ConfigFile)
}

func NewFeatureSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	if _, _, err := resolveConfig(ctx.FeatureName, ctx.Config); err != nil {
		return nil, err
	}
	return FeatureSnapshotGatherer{}, nil
}

func FeatureSmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return featurekit.SmokeSpec{
		WantMetrics: []string{metricExampleValue(ctx.FeatureName) + " 1"},
	}
}

func resolveConfig(featureName string, config Config) (string, bool, error) {
	var fileConfig configFile
	return loadFeatureConfigFile(featureName, config.ConfigFile, &fileConfig)
}
