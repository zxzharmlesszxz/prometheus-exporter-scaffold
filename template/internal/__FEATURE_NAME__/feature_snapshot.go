package __FEATURE_NAME__

import (
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

type FeatureSnapshotSpec struct {
	statusFunc func(Snapshot) framework.SnapshotStatus
}

func NewFeatureSnapshotSpec(statusFunc func(Snapshot) framework.SnapshotStatus) FeatureSnapshotSpec {
	return FeatureSnapshotSpec{statusFunc: statusFunc}
}

func (s FeatureSnapshotSpec) Status(snapshot Snapshot) framework.SnapshotStatus {
	if s.statusFunc == nil {
		return framework.SnapshotStatus{}
	}
	return s.statusFunc(snapshot)
}

func FeatureSnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}
