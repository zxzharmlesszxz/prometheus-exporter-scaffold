package exporter

import (
	"testing"
	"time"

	"github.com/alecthomas/kingpin/v2"
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
	if feature.refreshInterval != 30*time.Second {
		t.Fatalf("refreshInterval = %v, want %v", feature.refreshInterval, 30*time.Second)
	}
}
