package exporter

import (
	"io"
	"log/slog"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func testFeatureContext() framework.FeatureContext {
	return framework.FeatureContext{
		Logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		ExporterName: "__PROJECT_NAME__",
		Namespace:    defaultMetricNamespace,
	}
}
