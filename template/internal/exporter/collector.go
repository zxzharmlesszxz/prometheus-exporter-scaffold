package exporter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Snapshot struct {
	AttemptTime time.Time
	Success     bool
	Value       float64
	Err         error
}

type Snapshotter interface {
	Snapshot(context.Context, time.Time) Snapshot
}

type SnapshotGatherer struct{}

func (SnapshotGatherer) Snapshot(_ context.Context, now time.Time) Snapshot {
	return Snapshot{
		AttemptTime: now,
		Success:     true,
		Value:       1,
	}
}

type Collector struct {
	namespace       string
	logger          *slog.Logger
	snapshotter     Snapshotter
	refreshInterval time.Duration
	now             func() time.Time

	mu                       sync.Mutex
	initialized              bool
	backgroundStarted        bool
	snapshot                 Snapshot
	lastSuccessfulCollection time.Time

	exampleValueDesc             *prometheus.Desc
	lastCollectionSuccessDesc    *prometheus.Desc
	lastCollectionTimestampDesc  *prometheus.Desc
	lastSuccessfulCollectionDesc *prometheus.Desc
}

func NewCollector(namespace string, logger *slog.Logger, snapshotter Snapshotter, refreshInterval time.Duration) *Collector {
	if namespace == "" {
		namespace = "__METRIC_NAMESPACE__"
	}
	if logger == nil {
		logger = slog.Default()
	}
	if snapshotter == nil {
		snapshotter = SnapshotGatherer{}
	}
	if refreshInterval <= 0 {
		refreshInterval = defaultRefreshInterval
	}

	return &Collector{
		namespace:       namespace,
		logger:          logger,
		snapshotter:     snapshotter,
		refreshInterval: refreshInterval,
		now:             time.Now,

		exampleValueDesc: prometheus.NewDesc(
			"__FEATURE_NAME___example_value",
			"Example __FEATURE_NAME__ metric emitted by the generated exporter skeleton",
			nil,
			nil,
		),
		lastCollectionSuccessDesc: prometheus.NewDesc(
			namespace+"_last_collection_success",
			"Whether the last __FEATURE_NAME__ data collection succeeded",
			nil,
			nil,
		),
		lastCollectionTimestampDesc: prometheus.NewDesc(
			namespace+"_last_collection_timestamp_seconds",
			"Unix timestamp of the last __FEATURE_NAME__ data collection attempt",
			nil,
			nil,
		),
		lastSuccessfulCollectionDesc: prometheus.NewDesc(
			namespace+"_last_successful_collection_timestamp_seconds",
			"Unix timestamp of the last successful __FEATURE_NAME__ data collection",
			nil,
			nil,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.exampleValueDesc
	ch <- c.lastCollectionSuccessDesc
	ch <- c.lastCollectionTimestampDesc
	ch <- c.lastSuccessfulCollectionDesc
}

func (c *Collector) Start(ctx context.Context) {
	c.mu.Lock()
	if c.backgroundStarted {
		c.mu.Unlock()
		return
	}
	c.backgroundStarted = true
	c.mu.Unlock()

	go c.refreshLoop(ctx)
}

func (c *Collector) refreshLoop(ctx context.Context) {
	c.refresh(ctx, c.now())

	ticker := time.NewTicker(c.refreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.refresh(ctx, c.now())
		}
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	snapshot, lastSuccessful := c.currentSnapshot(c.now())

	ch <- prometheus.MustNewConstMetric(
		c.lastCollectionSuccessDesc,
		prometheus.GaugeValue,
		boolFloat(snapshot.Success),
	)
	ch <- prometheus.MustNewConstMetric(
		c.lastCollectionTimestampDesc,
		prometheus.GaugeValue,
		unixTimestamp(snapshot.AttemptTime),
	)
	ch <- prometheus.MustNewConstMetric(
		c.lastSuccessfulCollectionDesc,
		prometheus.GaugeValue,
		unixTimestamp(lastSuccessful),
	)

	if !snapshot.Success {
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.exampleValueDesc,
		prometheus.GaugeValue,
		snapshot.Value,
	)
}

func (c *Collector) currentSnapshot(now time.Time) (Snapshot, time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.backgroundStarted {
		return c.snapshot, c.lastSuccessfulCollection
	}
	if c.initialized && now.Sub(c.snapshot.AttemptTime) < c.refreshInterval {
		return c.snapshot, c.lastSuccessfulCollection
	}

	c.refreshLocked(context.Background(), now)
	return c.snapshot, c.lastSuccessfulCollection
}

func (c *Collector) refresh(ctx context.Context, now time.Time) {
	snapshot := c.snapshotter.Snapshot(ctx, now)
	if snapshot.Err != nil {
		c.logger.Error("__FEATURE_NAME__ data collection failed", "err", snapshot.Err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.storeSnapshot(snapshot)
}

func (c *Collector) refreshLocked(ctx context.Context, now time.Time) {
	snapshot := c.snapshotter.Snapshot(ctx, now)
	if snapshot.Err != nil {
		c.logger.Error("__FEATURE_NAME__ data collection failed", "err", snapshot.Err)
	}
	c.storeSnapshot(snapshot)
}

func (c *Collector) storeSnapshot(snapshot Snapshot) {
	if snapshot.Success {
		c.lastSuccessfulCollection = snapshot.AttemptTime
	}

	c.snapshot = snapshot
	c.initialized = true
}

func boolFloat(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

func unixTimestamp(value time.Time) float64 {
	if value.IsZero() {
		return 0
	}
	return float64(value.Unix())
}
