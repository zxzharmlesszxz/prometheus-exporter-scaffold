package exporter

import (
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter/exportertest"
)

func TestFeatureRegistersAndParsesFlags(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	app := kingpin.New("test", "")
	app.Terminate(func(int) {})
	feature.RegisterFlags(app)

	if _, err := app.Parse([]string{"--" + defaultFeatureName + ".refresh-interval=30s"}); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := feature.RuntimeConfig()
	if got := exportertest.RuntimeConfigValue(t, config, "refresh_interval"); got != 30*time.Second {
		t.Fatalf("refresh_interval = %v, want %v", got, 30*time.Second)
	}
}
