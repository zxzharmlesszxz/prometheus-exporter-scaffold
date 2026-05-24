package exporter

import (
	"testing"

	feature "__GO_MODULE__/internal/__FEATURE_NAME__"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestFeatureRuntimeConfigNormalizesValues(t *testing.T) {
	t.Parallel()

	exporterFeature := NewFeature()
	config := exporterFeature.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != feature.DefaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, feature.DefaultRefreshInterval)
	}
}
