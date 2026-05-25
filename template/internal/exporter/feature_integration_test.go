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

	info := ExporterInfo()
	wantMetrics := append([]string{info.Metrics.BuildInfo}, info.Smoke.WantMetrics...)

	handler := newTestHandler(t)
	body := waitForHandlerMetrics(t, handler, wantMetrics)
	if body == "" {
		t.Fatal("waitForHandlerMetrics() returned empty body")
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	feature := NewFeature()
	registry, err := framework.NewRegistry(
		defaultMetricNamespace,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		feature,
	)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v, want nil", err)
	}

	return framework.NewHandler(framework.HandlerOptions{
		Name:        defaultExporterName,
		Description: defaultExporterDescription,
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
