package __FEATURE_NAME__

import (
	"time"

	"github.com/alecthomas/kingpin/v2"
	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"
)

// Feature is the standard contract implemented by a concrete exporter feature.
type Feature interface {
	featurekit.FeatureContract[Config, Snapshot]
}

// FeatureExtension carries this exporter's feature-specific spec.
type FeatureExtension struct {
	featurekit.FeatureDefaults[Config, Snapshot]
	spec FeatureSpec
}

type FeatureSpec struct {
	refreshInterval time.Duration
	config          FeatureConfigSpec
	snapshot        FeatureSnapshotSpec
	snapshotter     FeatureSnapshotterSpec
	metrics         FeatureMetricsSpec
	smoke           FeatureSmokeSpec
}

type FeatureSmokeSpec struct {
	factory func(featurekit.SmokeContext[Config]) featurekit.SmokeSpec
}

func NewFeatureSmokeSpec(factory func(featurekit.SmokeContext[Config]) featurekit.SmokeSpec) FeatureSmokeSpec {
	return FeatureSmokeSpec{factory: factory}
}

func (s FeatureSmokeSpec) New(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	if s.factory == nil {
		return featurekit.SmokeSpec{}
	}
	return s.factory(ctx)
}

func (s FeatureSpec) DefaultRefreshInterval() time.Duration {
	return s.refreshInterval
}

func (s FeatureSpec) DefaultConfig() Config {
	return s.config.DefaultConfig()
}

func (s FeatureSpec) RegisterFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	s.config.RegisterFlags(app, ctx, config)
}

func (s FeatureSpec) ValidateConfig(config Config) error {
	return s.config.ValidateConfig(config)
}

func (s FeatureSpec) NewSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return s.snapshotter.New(ctx)
}

func (s FeatureSpec) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return s.snapshotter.DefaultSnapshotter()
}

func (s FeatureSpec) NewMetrics(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	return s.metrics.New(ctx)
}

func (s FeatureSpec) SnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return s.snapshot.Status(snapshot)
}

func (s FeatureSpec) RuntimeConfig(ctx featurekit.RuntimeConfigContext[Config]) []any {
	return s.config.RuntimeConfig(ctx)
}

func (s FeatureSpec) SmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return s.smoke.New(ctx)
}

var _ Feature = FeatureExtension{}

func NewFeatureExtension() FeatureExtension {
	return FeatureExtension{spec: NewFeatureSpec()}
}

func NewFeatureContract() Feature {
	return NewFeatureExtension()
}

func NewFeature(options featurekit.SpecOptions) *featurekit.Feature[Config, Snapshot] {
	return featurekit.NewFeature(featurekit.NewContractSnapshotFeatureSpec[Config, Snapshot](
		options,
		NewFeatureContract(),
	))
}

func (f FeatureExtension) DefaultRefreshInterval() time.Duration {
	return f.spec.DefaultRefreshInterval()
}

func (f FeatureExtension) DefaultConfig() Config {
	return f.spec.DefaultConfig()
}

func (f FeatureExtension) RegisterFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	f.spec.RegisterFlags(app, ctx, config)
}

func (f FeatureExtension) ValidateConfig(config Config) error {
	return f.spec.ValidateConfig(config)
}

func (f FeatureExtension) NewSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	return f.spec.NewSnapshotter(ctx)
}

func (f FeatureExtension) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return f.spec.DefaultSnapshotter()
}

func (f FeatureExtension) NewMetrics(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	return f.spec.NewMetrics(ctx)
}

func (f FeatureExtension) SnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	return f.spec.SnapshotStatus(snapshot)
}

func (f FeatureExtension) RuntimeConfig(ctx featurekit.RuntimeConfigContext[Config]) []any {
	return f.spec.RuntimeConfig(ctx)
}

func (f FeatureExtension) SmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	return f.spec.SmokeSpec(ctx)
}
