package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func (f *Feature) RegisterCollectors(ctx framework.FeatureContext, registry *prometheus.Registry) error {
	return f.feature.RegisterCollectors(ctx, registry)
}
