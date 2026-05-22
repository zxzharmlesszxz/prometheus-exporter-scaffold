# Context For Next Chat

This repository is `__PROJECT_NAME__`.

It was rendered from `prometheus-exporter-scaffold` and uses the framework:

```go
github.com/zxzharmlesszxz/prometheus-exporter-framework v0.1.4
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
- `exporter/exportertest/smoketest` for binary smoke tests
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
  - feature defaults, state, and constructor
- `internal/exporter/feature_flags.go`
  - feature flag registration
- `internal/exporter/feature_collectors.go`
  - feature collector registration
- `internal/exporter/runtime_config.go`
  - feature runtime config reporting
- `internal/exporter/main.go`
  - stable `Main()` framework bootstrap
- `internal/exporter/identity.go`
  - stable `FeatureName()` and `DefaultListenAddress()` methods
- `internal/exporter/listen.go`
  - rendered default listen address
- `internal/exporter/standard_metrics.go`
  - rendered feature name, metric namespace, build-info metric, and standard
    collection status metric constants
- `internal/exporter/metrics.go`
  - placeholder domain/example metric constants
- `internal/exporter/collector.go`
  - snapshot-backed placeholder collector
  - example metric descriptor
  - common collection status metrics through the framework
- `internal/exporter/collector_metrics.go`
  - placeholder metric description and emission methods
- `internal/exporter/collector_types.go`
  - placeholder collector, snapshot, and snapshot gatherer type declarations
- `internal/exporter/snapshot.go`
  - placeholder snapshot gathering plus snapshot status/error adapters
- `internal/exporter/*_test.go`
  - collector and feature tests
- `smoke/binary_test.go`
  - short `smoketest.Config`-based binary smoke test

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

`make check` also validates Prometheus and Docker Compose examples.

`make full-check` adds Docker smoke and release smoke.

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

- If `prometheus-exporter-framework v0.1.4` is not published yet, add a temporary
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

## Maintenance Rule

Keep this file updated whenever flags, metric namespace, domain packages,
framework version, Docker behavior, or verification state changes.
