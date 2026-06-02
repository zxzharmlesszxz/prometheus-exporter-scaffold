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
	RefreshInterval    time.Duration
	Config             Config
	RegisterFlagsFunc  func(*kingpin.Application, featurekit.FlagContext, *Config)
	ValidateConfigFunc func(Config) error
	NewSnapshotterFunc func(featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error)
	DefaultSnapshotter framework.Snapshotter[Snapshot]
	MetricsFunc        func(featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot]
	StatusFunc         func(Snapshot) framework.SnapshotStatus
	RuntimeConfigFunc  func(featurekit.RuntimeConfigContext[Config]) []any
	SmokeFunc          func(featurekit.SmokeContext[Config]) featurekit.SmokeSpec
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
	return f.spec.RefreshInterval
}

func (f FeatureExtension) DefaultConfig() Config {
	return f.spec.Config
}

func (f FeatureExtension) RegisterFlags(app *kingpin.Application, ctx featurekit.FlagContext, config *Config) {
	if f.spec.RegisterFlagsFunc != nil {
		f.spec.RegisterFlagsFunc(app, ctx, config)
	}
}

func (f FeatureExtension) ValidateConfig(config Config) error {
	if f.spec.ValidateConfigFunc == nil {
		return nil
	}
	return f.spec.ValidateConfigFunc(config)
}

func (f FeatureExtension) NewSnapshotter(ctx featurekit.CollectorContext[Config]) (framework.Snapshotter[Snapshot], error) {
	if f.spec.NewSnapshotterFunc == nil {
		return nil, nil
	}
	return f.spec.NewSnapshotterFunc(ctx)
}

func (f FeatureExtension) DefaultSnapshotter() framework.Snapshotter[Snapshot] {
	return f.spec.DefaultSnapshotter
}

func (f FeatureExtension) NewMetrics(ctx featurekit.SnapshotMetricsContext[Snapshot]) featurekit.SnapshotMetrics[Snapshot] {
	if f.spec.MetricsFunc == nil {
		return nil
	}
	return f.spec.MetricsFunc(ctx)
}

func (f FeatureExtension) SnapshotStatus(snapshot Snapshot) framework.SnapshotStatus {
	if f.spec.StatusFunc == nil {
		return framework.SnapshotStatus{}
	}
	return f.spec.StatusFunc(snapshot)
}

func (f FeatureExtension) RuntimeConfig(ctx featurekit.RuntimeConfigContext[Config]) []any {
	if f.spec.RuntimeConfigFunc == nil {
		return nil
	}
	return f.spec.RuntimeConfigFunc(ctx)
}

func (f FeatureExtension) SmokeSpec(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	if f.spec.SmokeFunc == nil {
		return featurekit.SmokeSpec{}
	}
	return f.spec.SmokeFunc(ctx)
}
