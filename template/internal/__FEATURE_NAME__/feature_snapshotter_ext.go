package __FEATURE_NAME__

import (
	"context"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type FeatureSnapshotGatherer struct{}

func NewDefaultSnapshotter() FeatureSnapshotGatherer {
	return FeatureSnapshotGatherer{}
}

func NewFeatureSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	if _, _, _, err := ResolveFeatureConfig(ctx.FeatureName, ctx.Config); err != nil {
		return nil, err
	}
	return FeatureSnapshotGatherer{}, nil
}

func FeatureSnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}

func (FeatureSnapshotGatherer) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}
