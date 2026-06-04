package __FEATURE_NAME__

import (
	"context"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

func NewDefaultSnapshotEngine() featurekit.SnapshotEngine[Snapshot] {
	engine, err := newSnapshotEngine(NewDefaultConfig())
	if err != nil {
		panic(err)
	}
	return engine
}

func NewSnapshotEngine(ctx featurekit.CollectorContext[Config]) (featurekit.SnapshotEngine[Snapshot], error) {
	config, _, _, err := ResolveFeatureConfig(ctx.FeatureName, ctx.Config)
	if err != nil {
		return nil, err
	}
	return newSnapshotEngine(config)
}

func FeatureSnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}

func newSnapshotEngine(_ Config) (featurekit.SnapshotEngine[Snapshot], error) {
	return defaultSnapshotEngine{}, nil
}

type defaultSnapshotEngine struct{}

func (e defaultSnapshotEngine) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}
