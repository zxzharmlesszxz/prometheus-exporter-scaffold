package __FEATURE_NAME__

import (
	"context"
	"time"

	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type SnapshotGathererConfig struct{}

type SnapshotGatherer struct{}

func DefaultSnapshotGathererConfig() SnapshotGathererConfig {
	return SnapshotGathererConfig{}
}

func NewSnapshotGathererConfig(ctx featurekit.CollectorContext[Config]) (SnapshotGathererConfig, error) {
	if _, _, _, err := ResolveFeatureConfig(ctx.FeatureName, ctx.Config); err != nil {
		return SnapshotGathererConfig{}, err
	}
	return DefaultSnapshotGathererConfig(), nil
}

func NewSnapshotGatherer(_ SnapshotGathererConfig) SnapshotGatherer {
	return SnapshotGatherer{}
}

func (SnapshotGatherer) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}
