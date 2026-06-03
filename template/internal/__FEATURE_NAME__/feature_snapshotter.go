package __FEATURE_NAME__

import (
	"context"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type SnapshotEngine interface {
	Snapshot(context.Context, time.Time) Snapshot
}

type SnapshotGatherer struct {
	engine SnapshotEngine
}

func NewDefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return NewSnapshotGatherer(NewDefaultSnapshotEngine())
}

func NewFeatureSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	engine, err := NewSnapshotEngine(ctx)
	if err != nil {
		return nil, err
	}
	return NewSnapshotGatherer(engine), nil
}

func FeatureSnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}

func NewSnapshotGatherer(engine SnapshotEngine) SnapshotGatherer {
	return SnapshotGatherer{
		engine: engine,
	}
}

func (g SnapshotGatherer) Snapshot(ctx context.Context, now time.Time) Snapshot {
	engine := g.engine
	if engine == nil {
		engine = NewDefaultSnapshotEngine()
	}
	return engine.Snapshot(ctx, now)
}
