package __FEATURE_NAME__

import (
	"context"
	"time"

	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

func NewDefaultSnapshotEngine() SnapshotEngine {
	return defaultSnapshotEngine{}
}

func NewSnapshotEngine(ctx featurekit.CollectorContext[Config]) (SnapshotEngine, error) {
	if _, _, _, err := ResolveFeatureConfig(ctx.FeatureName, ctx.Config); err != nil {
		return nil, err
	}
	return NewDefaultSnapshotEngine(), nil
}

type defaultSnapshotEngine struct{}

func (e defaultSnapshotEngine) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}
