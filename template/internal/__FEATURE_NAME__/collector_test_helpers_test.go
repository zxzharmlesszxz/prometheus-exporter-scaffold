package __FEATURE_NAME__

import (
	"context"
	"sync/atomic"
	"time"
)

const (
	testFeatureName      = "__FEATURE_NAME__"
	testMetricNamespace  = "__METRIC_NAMESPACE__"
	testRefreshInterval  = time.Minute
	testLastSuccess      = testMetricNamespace + "_last_collection_success"
	testLastTimestamp    = testMetricNamespace + "_last_collection_timestamp_seconds"
	testLastSuccessfulTS = testMetricNamespace + "_last_successful_collection_timestamp_seconds"
)

type fakeSnapshotter struct {
	snapshot atomic.Value
}

func newFakeSnapshotter(snapshot Snapshot) *fakeSnapshotter {
	s := &fakeSnapshotter{}
	s.snapshot.Store(snapshot)
	return s
}

func (s *fakeSnapshotter) Snapshot(context.Context, time.Time) Snapshot {
	return s.snapshot.Load().(Snapshot)
}

func (s *fakeSnapshotter) set(snapshot Snapshot) {
	s.snapshot.Store(snapshot)
}
