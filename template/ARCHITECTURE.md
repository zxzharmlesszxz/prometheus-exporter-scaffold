# Architecture

`__PROJECT_NAME__` is a thin concrete exporter built on `prometheus-exporter-framework`.

## Package Layout

- `cmd`
  Minimal process entrypoint.
- `internal/exporter`
  Thin adapter that creates the domain feature and delegates bootstrap metadata to the framework.
- `internal/__FEATURE_NAME__`
  Feature spec, domain config, snapshot type, metric descriptors, and typed snapshot-to-metrics adapter.
- `smoke`
  Binary smoke tests that build the real executable and verify CLI, HTTP, and metric behavior.

## Data Flow

1. `cmd/main.go` delegates to `internal/exporter.Main()`, which runs `framework.MainFromInjectedProject(...)`.
2. `internal/exporter` creates the concrete feature from `internal/__FEATURE_NAME__` using framework-injected feature metadata.
3. Framework `featurekit.Feature` registers common flags such as `--__FEATURE_NAME__.refresh-interval` and delegates domain flags to the typed spec.
4. Framework `featurekit.Feature` builds a typed snapshotter and collector from the spec, then registers and starts the collector.
5. `framework.SnapshotCollector` refreshes data in a background worker every `--__FEATURE_NAME__.refresh-interval`; scrapes read the latest completed snapshot.
6. The collector exports domain metrics and collection health metrics.

## Failure Semantics

If refresh fails before the first successful snapshot, the exporter exposes collection health metrics, but no business metrics.

If the latest refresh fails, the exporter exposes collection health metrics, but no business metrics, and sets:

- `__METRIC_NAMESPACE___last_collection_success = 0`

The `/healthz` endpoint remains `200 OK` while the process is alive even if the latest collection failed.
