package __FEATURE_NAME__

import (
	"context"
	"log/slog"
	"time"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func (SnapshotGatherer) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}

func snapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return framework.SnapshotStatus{
		AttemptTime: snapshot.AttemptTime,
		Success:     snapshot.Success,
	}
}

func (m *Metrics) LogSnapshotError(logger *slog.Logger, snapshot Snapshot) {
	if snapshot.Err != nil {
		logger.Error(m.featureName+" data collection failed", "err", snapshot.Err)
	}
}
