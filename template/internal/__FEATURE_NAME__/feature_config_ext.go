package __FEATURE_NAME__

import (
	"time"

	"github.com/alecthomas/kingpin/v2"
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

func FeatureConfigFile(config *Config) *string {
	return &config.ConfigFile
}

func RegisterFeatureConfigFlags(*kingpin.Application, featurekit.FlagContext, *Config) {}

func ValidateFeatureConfig(Config) error {
	return nil
}

func FeatureRuntimeConfigEntries(featurekit.RuntimeConfigContext[Config], Config) []any {
	return nil
}

func ResolveFeatureConfig(featureName string, config Config) (Config, string, bool, error) {
	var fileConfig configFile
	configFile, loaded, err := loadFeatureConfigFile(featureName, config.ConfigFile, &fileConfig)
	return config, configFile, loaded, err
}
