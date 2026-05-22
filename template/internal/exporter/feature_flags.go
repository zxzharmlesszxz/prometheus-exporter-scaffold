package exporter

import "github.com/alecthomas/kingpin/v2"

func (f *Feature) RegisterFlags(app *kingpin.Application) {
	app.Flag(
		defaultFeatureName+".refresh-interval", "How often exporter refreshes "+defaultFeatureName+" data",
	).Default(defaultRefreshInterval.String()).DurationVar(&f.refreshInterval)
}
