# Architecture

`__PROJECT_NAME__` is a thin concrete exporter built on `prometheus-template-exporter`.

## Package Layout

- `cmd`
  Minimal process entrypoint.
- `internal/exporter`
  Template feature, flags, domain snapshot type, metric descriptors, and typed snapshot-to-metrics adapter.
- `smoke`
  Binary smoke tests that build the real executable and verify CLI, HTTP, and metric behavior.

## Data Flow

1. `cmd/main.go` delegates to `internal/exporter.Main()`, which runs `template.MainFromProject(...)`.
2. The feature registers `--__FEATURE_NAME__.refresh-interval`.
3. The feature registers one collector with the template registry.
4. `template.SnapshotCollector` refreshes data in a background worker every `--__FEATURE_NAME__.refresh-interval`; scrapes read the latest completed snapshot.
5. The collector exports domain metrics and collection health metrics.

## Failure Semantics

If refresh fails before the first successful snapshot, the exporter exposes collection health metrics, but no business metrics.

If the latest refresh fails, the exporter exposes collection health metrics, but no business metrics, and sets:

- `__METRIC_NAMESPACE___last_collection_success = 0`

The `/healthz` endpoint remains `200 OK` while the process is alive even if the latest collection failed.
