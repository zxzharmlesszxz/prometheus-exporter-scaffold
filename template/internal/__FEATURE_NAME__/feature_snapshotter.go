package __FEATURE_NAME__

import (
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

func NewDefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return NewSnapshotGatherer(DefaultSnapshotGathererConfig())
}

func NewFeatureSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	gathererConfig, err := NewSnapshotGathererConfig(ctx)
	if err != nil {
		return nil, err
	}
	return NewSnapshotGatherer(gathererConfig), nil
}

func FeatureSnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}
