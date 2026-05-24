# prometheus-exporter-scaffold

Scaffold repository for creating concrete Prometheus exporters from
`prometheus-exporter-framework`.

This repository owns generated-repository shape:

- project layout
- placeholder exporter feature
- typed snapshot collector wiring
- shared exporter test helper usage
- example Prometheus, Grafana, and Docker Compose files
- GitHub Actions and GitLab CI starter workflows
- Dependabot starter configuration
- rendering script

The framework itself lives in
`github.com/zxzharmlesszxz/prometheus-exporter-framework`. Keep exporter runtime
behavior in the framework and keep generated-project boilerplate here.

## Render A New Exporter

```bash
scripts/render.sh \
  --project-name prometheus-demo-exporter \
  --module github.com/example/prometheus-demo-exporter \
  --description "Prometheus Demo Exporter" \
  --feature-name demo \
  --namespace demo_exporter \
  --port 9888 \
  --target-dir /tmp/prometheus-demo-exporter
```

Then validate the generated repository:

```bash
cd /tmp/prometheus-demo-exporter
go mod tidy
make go-check
make check
```

`--feature-name`, `--namespace`, `--description`, `--module`, and `--port`
have defaults, but passing them explicitly keeps the generated repository
predictable.

The generated `cmd/main.go` is intentionally stable. Project metadata is
injected by Makefile linker flags from `Makefile.mk`, while the concrete feature
package owns domain behavior.

## Framework Version

`framework.version` and `template/go.mod` track the
`prometheus-exporter-framework` version used by newly generated exporters.

When a new framework tag is published, the framework release workflow opens a
pull request here to update the scaffold and verify a rendered exporter.

This repository's own CI also renders a demo exporter and runs its Go-only
checks, so scaffold pull requests validate the generated code path directly.

## Update An Existing Exporter

Existing exporters are not coupled to this repository after rendering. To check
or sync scaffold-owned files against the current template, run:

```bash
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter
```

To update the default managed files:

```bash
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --sync
```

The default managed set is intentionally conservative: CI files, ignore files,
`cmd/main.go`, Dependabot config, and stable scaffold-owned Go wiring under `internal/exporter`.
Concrete exporters keep domain logic in their rendered feature package, normally
`internal/<feature-name>`, so inspect those files separately instead of blindly
syncing them:

```bash
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file Makefile
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file Dockerfile
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/exporter/feature_flags_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/exporter/feature_collectors_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/exporter/runtime_config_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/exporter/identity_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/exporter/feature_test_helpers_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/exporter/feature_integration_test_helpers_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/exporter.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/metrics.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/collector_types.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/collector_metrics.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/collector_test_helpers_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/collector_snapshot_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/collector_refresh_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/collector_defaults_test.go
scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --file internal/demo/snapshot.go
```

Older exporters may still define `Main()`, `FeatureName()`, or
`DefaultListenAddress()` inside `internal/exporter/feature.go`, and
`defaultListenAddress` may still live there as well. Remove those definitions
once when adopting the split scaffold Go files; the drift script will report
`LEGACY` instead of syncing duplicate definitions. It reports the same guard for
standard metric constants that still live in `metrics.go`, and for collector
types, collector metric methods, and snapshot helpers that still live in
`internal/exporter/collector.go`. Collector test helpers get the same guard
while they still live in `internal/exporter/collector_test.go`; collector tests
get it while they still live in `collector_test.go`; domain feature methods and
test helpers get it while they still live in `internal/exporter/feature.go`,
`feature_test.go`, or `feature_integration_test.go`.
