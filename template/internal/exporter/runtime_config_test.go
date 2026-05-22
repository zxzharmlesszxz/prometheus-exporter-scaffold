package exporter

import (
	"testing"
	"time"

	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestFeatureRuntimeConfigNormalizesValues(t *testing.T) {
	t.Parallel()

	feature := &Feature{refreshInterval: -time.Second}
	config := feature.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != defaultRefreshInterval {
		t.Fatalf("refresh_interval = %v, want %v", got, defaultRefreshInterval)
	}
}
