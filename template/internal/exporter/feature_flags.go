package exporter

import "github.com/alecthomas/kingpin/v2"

func (f *Feature) RegisterFlags(app *kingpin.Application) {
	f.feature.RegisterFlags(app)
}
