package __FEATURE_NAME__

import "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/featurekit"

type FeatureSmokeFactory func(featurekit.SmokeContext[Config]) featurekit.SmokeSpec

type FeatureSmokeSpec struct {
	factory FeatureSmokeFactory
}

func NewFeatureSmokeSpec(factory FeatureSmokeFactory) FeatureSmokeSpec {
	return FeatureSmokeSpec{factory: factory}
}

func (s FeatureSmokeSpec) New(ctx featurekit.SmokeContext[Config]) featurekit.SmokeSpec {
	if s.factory == nil {
		return featurekit.SmokeSpec{}
	}
	return s.factory(ctx)
}
