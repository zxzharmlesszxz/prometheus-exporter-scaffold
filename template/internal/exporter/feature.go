package exporter

import "time"

const (
	defaultRefreshInterval = time.Minute
)

type Feature struct {
	refreshInterval time.Duration
}

func NewFeature() *Feature {
	return &Feature{
		refreshInterval: defaultRefreshInterval,
	}
}
