package __FEATURE_NAME__

import (
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func FeatureSnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}
