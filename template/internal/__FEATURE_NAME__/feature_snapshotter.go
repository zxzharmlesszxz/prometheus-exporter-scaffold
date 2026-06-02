package __FEATURE_NAME__

import (
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

type FeatureSnapshotterFactory func(featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error)

type FeatureSnapshotterSpec struct {
	factory            FeatureSnapshotterFactory
	defaultSnapshotter framework.Snapshotter[Snapshot]
}

func NewFeatureSnapshotterSpec(factory FeatureSnapshotterFactory, defaultSnapshotter framework.Snapshotter[Snapshot]) FeatureSnapshotterSpec {
	return FeatureSnapshotterSpec{
		factory:            factory,
		defaultSnapshotter: defaultSnapshotter,
	}
}

func (s FeatureSnapshotterSpec) New(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	if s.factory == nil {
		return nil, nil
	}
	return s.factory(ctx)
}

func (s FeatureSnapshotterSpec) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return s.defaultSnapshotter
}
