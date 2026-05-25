package __FEATURE_NAME__

import (
	"time"
)

type Snapshot struct {
	AttemptTime time.Time
	Success     bool
	Value       float64
	Err         error
}

type SnapshotGatherer struct{}
