package __FEATURE_NAME__

func metricExampleValue(featureName string) string {
	return featureName + "_example_value"
}

const (
	metricLastCollectionSuccess                    = defaultMetricNamespace + "_last_collection_success"
	metricLastCollectionTimestampSeconds           = defaultMetricNamespace + "_last_collection_timestamp_seconds"
	metricLastSuccessfulCollectionTimestampSeconds = defaultMetricNamespace + "_last_successful_collection_timestamp_seconds"
)
