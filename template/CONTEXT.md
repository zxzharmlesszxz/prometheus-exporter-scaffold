# Context For Next Chat

This repository is `__PROJECT_NAME__`.

It was rendered from `prometheus-exporter-scaffold` and uses the framework:

```go
github.com/zxzharmlesszxz/prometheus-exporter-framework v0.2.0
```

Local Go tooling in the original workspace is expected at:

```bash
PATH=/Users/mort/sdk/go1.26.3/bin:$PATH
```

## Framework Context

`prometheus-exporter-framework` provides the reusable exporter shell:

- CLI bootstrap through `framework.MainFromProject`
- dynamic executable-name usage in `--help`
- standard Prometheus/exporter-toolkit flags
- `promslog` logging
- `/metrics`, `/healthz`, landing page, optional pprof
- custom Prometheus registry
- build info, Go runtime, and process collectors
- typed snapshot cache and background refresh collector helper
- `exporter/exportertest` helpers for collector tests
- `exporter/exportertest/smoketest` for binary smoke tests, including
  `smoketest.Config.BinaryPath` for Makefile-built binaries
- version metadata hydration from linker flags or Go build info

## Current Exporter State

Entrypoint:

```text
cmd/main.go
```

Generated entrypoint:

```go
package main

import "__GO_MODULE__/internal/exporter"

func main() {
	exporter.Main()
}
```

`internal/exporter.Main()` calls:

```go
framework.MainFromProject(NewFeature())
```

Feature name:

```text
__FEATURE_NAME__
```

Metric namespace:

```text
__METRIC_NAMESPACE__
```

Default listen address:

```text
:__DEFAULT_PORT__
```

Default refresh interval:

```text
1m
```

## Generated Domain Skeleton

- `internal/exporter/feature.go`
  - stable adapter constructor delegating to `internal/__FEATURE_NAME__`
- `internal/exporter/feature_flags.go`
  - stable adapter flag delegation
- `internal/exporter/feature_collectors.go`
  - stable adapter collector registration delegation
- `internal/exporter/runtime_config.go`
  - stable adapter runtime config delegation
- `internal/exporter/main.go`
  - stable `Main()` framework bootstrap
- `internal/exporter/identity.go`
  - stable `FeatureName()` and `DefaultListenAddress()` methods
- `internal/exporter/defaults.go`
  - Makefile-injected linker vars for exporter name, description, feature name,
    metric namespace, and default listen address
- `internal/exporter/standard_metrics.go`
  - build-info metric and standard collection status metric constants
- `internal/exporter/info.go`
  - stable exporter metadata and binary smoke configuration type
- `internal/exporter/info_test.go`
  - stable checks for common exporter metadata and smoke configuration
- `internal/exporter/feature_flags_test.go`,
  `feature_collectors_test.go`, `runtime_config_test.go`, and
  `identity_test.go`
  - stable adapter tests by concern
- `internal/exporter/feature_test_helpers_test.go` and
  `feature_integration_test_helpers_test.go`
  - stable adapter test helpers
- `internal/__FEATURE_NAME__/exporter.go`
  - placeholder domain flags, runtime config, and collector registration
- `internal/__FEATURE_NAME__/metrics.go`
  - placeholder domain/example metric constants
- `internal/__FEATURE_NAME__/smoke.go`
  - placeholder domain smoke additions consumed by `internal/exporter/info.go`
- `internal/__FEATURE_NAME__/collector.go`
  - snapshot-backed placeholder collector
  - example metric descriptor
  - common collection status metrics through the framework
- `internal/__FEATURE_NAME__/collector_metrics.go`
  - placeholder metric description and emission methods
- `internal/__FEATURE_NAME__/collector_types.go`
  - placeholder collector, snapshot, and snapshot gatherer type declarations
- `internal/__FEATURE_NAME__/snapshot.go`
  - placeholder snapshot gathering plus snapshot status/error adapters
- `internal/__FEATURE_NAME__/*_test.go`
  - collector and placeholder domain exporter tests
  - split collector tests by concern in `collector_snapshot_test.go`,
    `collector_refresh_test.go`, and `collector_defaults_test.go`
  - split collector test helper in `collector_test_helpers_test.go`
  - placeholder domain exporter tests in `exporter_test.go`
- `smoke/binary_test.go`
  - short `smoketest.Config`-based binary smoke test
  - imports `internal/exporter` and passes `ExporterInfo()` into the framework
    smoke helper

When turning this into a real exporter, replace placeholder domain logic and
examples while keeping the stable framework wiring.

## Verification Targets

```bash
make help
make go-check
make check
make docker-smoke
make full-check
```

`make go-check` runs formatting, vet, staticcheck, coverage threshold, binary
smoke, and race tests.

Use Make targets for Go builds/tests that import `internal/exporter`; raw
`go run ./cmd`, `go build ./cmd`, and `go test ./...` do not inject exporter
metadata and are intentionally unsupported.

`make check` also validates Prometheus and Docker Compose examples.

`make full-check` adds Docker smoke and release smoke.

## Makefile Shape

Common override-able variables live in:

```text
exporter.mk
```

`Makefile` includes `exporter.mk` and keeps target logic only. Existing
exporters may customize target bodies, especially Docker smoke checks, while
keeping `exporter.mk` scaffold-managed.

`PROJECT_NAME` is rendered to `__PROJECT_NAME__`; do not derive it from the
temporary render directory. `SMOKE_BINARY` defaults to `$(DIST_DIR)/$(PROJECT_NAME)`
so CLI usage based on executable basename still matches `ExporterInfo().Name`.

`exporter.mk` owns metadata linker flags. `LDFLAGS` and `SMOKE_LDFLAGS` inject:

```text
internal/exporter.defaultExporterName
internal/exporter.defaultExporterDescription
internal/exporter.defaultFeatureName
internal/exporter.defaultMetricNamespace
internal/exporter.defaultListenAddress
```

Docker smoke has one exporter-specific assertion controlled by:

```make
DOCKER_SMOKE_METRIC ?= $(FEATURE_NAME)_example_value 1
```

The `docker-smoke-image` target greps for that value after standard health,
build-info, and collection-success checks.

## Docker Notes

The runtime image copies the built project binary to:

```text
/usr/local/bin/exporter
```

So Docker `--help` shows:

```text
usage: exporter [<flags>]
```

even though local/release binaries use the project executable file name.

## Known Pending Work From Scaffold

- If `prometheus-exporter-framework v0.2.0` is not published yet, add a temporary
  local replace before running Go checks:

  ```go
  replace github.com/zxzharmlesszxz/prometheus-exporter-framework => ../prometheus-exporter-framework
  ```

- Remove temporary local replaces after the framework tag is published.
- Decide whether Dockerfile should move from `golang:1.26` +
  `debian:bookworm-slim` to Alpine images.

## Latest Verification

On 2026-05-21 the generated `Makefile` Docker build command was checked with:

```bash
make -n docker-build
```

The previous line-continuation issue in Docker targets is fixed.

On 2026-05-22 the scaffold template moved placeholder collector and snapshot
logic out of `internal/exporter` into `internal/__FEATURE_NAME__`. After render,
that path becomes the concrete domain package, for example `internal/demo` or
`internal/domain`. `internal/exporter` should stay a stable framework adapter
that delegates feature flags, collector registration, and runtime config to the
domain package.

On 2026-05-22 `internal/exporter/feature_integration_test.go` was changed to
derive expected metrics from `ExporterInfo()`: build-info comes from
`info.Metrics.BuildInfo`, and exporter-specific expectations come from
`info.Smoke.WantMetrics`. The placeholder example metric is added by
`internal/__FEATURE_NAME__/smoke.go`, so real exporters can keep those domain
smoke variables empty and keep the integration test scaffold-identical.

On 2026-05-22 exporter metadata moved to Makefile-only linker injection.
`internal/exporter/defaults.go` has empty string vars and fails fast during
package init when a supported Make target does not inject them. `make test`,
`make coverage`, `make test-race`, and `make smoke` pass `-ldflags`.
`make smoke` builds `$(SMOKE_BINARY)` first, then passes it through
`EXPORTER_SMOKE_BINARY` to `smoketest.Config.BinaryPath`.

A fresh rendered demo exporter was verified with a temporary local framework
replace: placeholder scan, `make go-check COVERAGE_THRESHOLD=90.0`, and default
scaffold drift passed with total coverage `92.7%`.

## Maintenance Rule

Keep this file updated whenever flags, metric namespace, domain packages,
framework version, Docker behavior, or verification state changes.
