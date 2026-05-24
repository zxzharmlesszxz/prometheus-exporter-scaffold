package exporter

import feature "__GO_MODULE__/internal/__FEATURE_NAME__"

type Exporter struct {
	Name                 string
	Description          string
	FeatureName          string
	MetricNamespace      string
	DefaultListenAddress string
	Metrics              MetricInfo
	Smoke                SmokeInfo
}

type MetricInfo struct {
	BuildInfo                                string
	LastCollectionSuccess                    string
	LastCollectionTimestampSeconds           string
	LastSuccessfulCollectionTimestampSeconds string
}

type SmokeInfo struct {
	ForbiddenUsageNames []string
	RenamedExecutable   string
	ServerArgs          []string
	WantMetrics         []string
	RejectMetrics       []string
}

func ExporterInfo() Exporter {
	return defaultExporterInfo()
}

func defaultExporterInfo() Exporter {
	return Exporter{
		Name:                 defaultExporterName,
		Description:          defaultExporterDescription,
		FeatureName:          defaultFeatureName,
		MetricNamespace:      defaultMetricNamespace,
		DefaultListenAddress: defaultListenAddress,
		Metrics: MetricInfo{
			BuildInfo:                                metricBuildInfo,
			LastCollectionSuccess:                    metricLastCollectionSuccess,
			LastCollectionTimestampSeconds:           metricLastCollectionTimestampSeconds,
			LastSuccessfulCollectionTimestampSeconds: metricLastSuccessfulCollectionTimestampSeconds,
		},
		Smoke: SmokeInfo{
			ForbiddenUsageNames: []string{defaultMetricNamespace},
			RenamedExecutable:   "renamed-" + defaultFeatureName + "-exporter",
			ServerArgs:          append([]string{"--" + defaultFeatureName + ".refresh-interval=100ms"}, feature.SmokeServerArgs...),
			WantMetrics:         append([]string{metricLastCollectionSuccess + " 1"}, feature.SmokeWantMetrics...),
			RejectMetrics:       append([]string{metricLastCollectionSuccess + " 0"}, feature.SmokeRejectMetrics...),
		},
	}
}
