package exporter

import "testing"

func TestExporterInfo(t *testing.T) {
	t.Parallel()

	info := ExporterInfo()
	if info.Name != defaultExporterName {
		t.Fatalf("Name = %q, want %q", info.Name, defaultExporterName)
	}
	if info.Description != defaultExporterDescription {
		t.Fatalf("Description = %q, want %q", info.Description, defaultExporterDescription)
	}
	if info.FeatureName != defaultFeatureName {
		t.Fatalf("FeatureName = %q, want %q", info.FeatureName, defaultFeatureName)
	}
	if info.MetricNamespace != defaultMetricNamespace {
		t.Fatalf("MetricNamespace = %q, want %q", info.MetricNamespace, defaultMetricNamespace)
	}
	if info.DefaultListenAddress != defaultListenAddress {
		t.Fatalf("DefaultListenAddress = %q, want %q", info.DefaultListenAddress, defaultListenAddress)
	}
	if info.Metrics.BuildInfo != metricBuildInfo {
		t.Fatalf("Metrics.BuildInfo = %q, want %q", info.Metrics.BuildInfo, metricBuildInfo)
	}
	if info.Metrics.LastCollectionSuccess != metricLastCollectionSuccess {
		t.Fatalf("Metrics.LastCollectionSuccess = %q, want %q", info.Metrics.LastCollectionSuccess, metricLastCollectionSuccess)
	}
	if !hasString(info.Smoke.ForbiddenUsageNames, defaultMetricNamespace) {
		t.Fatalf("Smoke.ForbiddenUsageNames = %v, want %q", info.Smoke.ForbiddenUsageNames, defaultMetricNamespace)
	}
	if info.Smoke.RenamedExecutable != "renamed-"+defaultFeatureName+"-exporter" {
		t.Fatalf("Smoke.RenamedExecutable = %q", info.Smoke.RenamedExecutable)
	}
	if !hasString(info.Smoke.ServerArgs, "--"+defaultFeatureName+".refresh-interval=100ms") {
		t.Fatalf("Smoke.ServerArgs = %v", info.Smoke.ServerArgs)
	}
	if !hasString(info.Smoke.WantMetrics, metricLastCollectionSuccess+" 1") {
		t.Fatalf("Smoke.WantMetrics = %v", info.Smoke.WantMetrics)
	}
	if !hasString(info.Smoke.RejectMetrics, metricLastCollectionSuccess+" 0") {
		t.Fatalf("Smoke.RejectMetrics = %v", info.Smoke.RejectMetrics)
	}
}

func hasString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
