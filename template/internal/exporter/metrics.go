package exporter

const (
	defaultFeatureName     = "__FEATURE_NAME__"
	defaultMetricNamespace = "__METRIC_NAMESPACE__"

	metricBuildInfo                                = defaultMetricNamespace + "_build_info"
	metricExampleValue                             = defaultFeatureName + "_example_value"
	metricLastCollectionSuccess                    = defaultMetricNamespace + "_last_collection_success"
	metricLastCollectionTimestampSeconds           = defaultMetricNamespace + "_last_collection_timestamp_seconds"
	metricLastSuccessfulCollectionTimestampSeconds = defaultMetricNamespace + "_last_successful_collection_timestamp_seconds"
)
