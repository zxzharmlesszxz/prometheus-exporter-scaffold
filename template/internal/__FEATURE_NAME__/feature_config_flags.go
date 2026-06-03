package __FEATURE_NAME__

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

func RegisterFeatureConfigFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	featurekit.RegisterFeatureConfigFlagSpecs(app, ctx, config, featureConfigFlagSpecs)
}
