package exporter

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func (f *Feature) RegisterCollectors(ctx framework.FeatureContext, registry *prometheus.Registry) error {
	collector := NewCollector(
		ctx.Namespace,
		ctx.Logger,
		SnapshotGatherer{},
		framework.NormalizeDuration(f.refreshInterval, defaultRefreshInterval),
	)
	if err := framework.RegisterAndStartCollectors(context.Background(), registry, collector); err != nil {
		return fmt.Errorf("register %s collector: %w", defaultFeatureName, err)
	}
	return nil
}
