package exporter

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	framework "github.com/zxzharmlesszxz/prometheus-exporter-framework/exporter"
)

func TestFeatureServesMetricsThroughTemplate(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)
	body := waitForHandlerMetrics(t, handler, []string{
		"__METRIC_NAMESPACE___build_info",
		"__METRIC_NAMESPACE___last_collection_success 1",
		"__FEATURE_NAME___example_value 1",
	})
	if body == "" {
		t.Fatal("waitForHandlerMetrics() returned empty body")
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	feature := NewFeature()
	registry, err := framework.NewRegistry(
		"__METRIC_NAMESPACE__",
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		feature,
	)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v, want nil", err)
	}

	return framework.NewHandler(framework.HandlerOptions{
		Name:        "__PROJECT_NAME__",
		Description: "__PROJECT_DESC__",
		MetricsPath: "/metrics",
		Registry:    registry,
	})
}

func waitForHandlerMetrics(t *testing.T, handler http.Handler, wants []string) string {
	t.Helper()

	for i := 0; i < 100; i++ {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("ServeHTTP() status = %d, want %d", rec.Code, http.StatusOK)
		}

		body := rec.Body.String()
		missing := ""
		for _, want := range wants {
			if !strings.Contains(body, want) {
				missing = want
				break
			}
		}
		if missing == "" {
			return body
		}
	}
	t.Fatalf("ServeHTTP() body did not contain metrics %v", wants)
	return ""
}
