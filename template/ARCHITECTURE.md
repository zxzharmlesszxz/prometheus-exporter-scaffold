# Architecture

`__PROJECT_NAME__` is a thin concrete exporter built on `prometheus-template-exporter`.

## Package Layout

- `cmd`
  Minimal process entrypoint.
- `internal/exporter`
  Template feature, flags, background refresh worker, snapshot cache, Prometheus collector, and metric descriptors.
- `smoke`
  Binary smoke tests that build the real executable and verify CLI, HTTP, and metric behavior.

## Data Flow

1. `cmd/main.go` delegates to `template.MainForProject(...)`.
2. The feature registers `--__FEATURE_NAME__.refresh-interval`.
3. The feature registers one collector with the template registry.
4. The collector refreshes data in a background worker every `--__FEATURE_NAME__.refresh-interval`; scrapes read the latest completed snapshot.
5. The collector exports domain metrics and collection health metrics.

## Failure Semantics

If refresh fails before the first successful snapshot, the exporter exposes collection health metrics, but no business metrics.

If refresh fails after at least one successful snapshot, keep serving the last successful business metrics and set:

- `__METRIC_NAMESPACE___last_collection_success = 0`

The `/healthz` endpoint remains `200 OK` while the process is alive even if the latest collection failed.
