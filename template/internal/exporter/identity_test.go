package exporter

import "testing"

func TestFeatureMetadata(t *testing.T) {
	t.Parallel()

	feature := NewFeature()
	if feature.FeatureName() != defaultFeatureName {
		t.Fatalf("FeatureName() = %q, want %q", feature.FeatureName(), defaultFeatureName)
	}
	if feature.DefaultListenAddress() != defaultListenAddress {
		t.Fatalf("DefaultListenAddress() = %q, want %q", feature.DefaultListenAddress(), defaultListenAddress)
	}
}
