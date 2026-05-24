package exporter

var (
	defaultExporterName        string
	defaultExporterDescription string
	defaultFeatureName         string
	defaultMetricNamespace     string
	defaultListenAddress       string
)

func init() {
	requireInjectedDefault("defaultExporterName", defaultExporterName)
	requireInjectedDefault("defaultExporterDescription", defaultExporterDescription)
	requireInjectedDefault("defaultFeatureName", defaultFeatureName)
	requireInjectedDefault("defaultMetricNamespace", defaultMetricNamespace)
	requireInjectedDefault("defaultListenAddress", defaultListenAddress)
	if defaultListenAddress[0] != ':' {
		panic("invalid Makefile-injected exporter metadata: defaultListenAddress must start with :")
	}
}

func requireInjectedDefault(name string, value string) {
	if value == "" {
		panic("missing Makefile-injected exporter metadata: " + name)
	}
}
