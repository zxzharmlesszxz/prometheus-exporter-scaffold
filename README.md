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
make new-exporter \
  PROJECT_NAME=prometheus-demo-exporter \
  GO_MODULE=github.com/example/prometheus-demo-exporter \
  PROJECT_DESC="Prometheus Demo Exporter" \
  FEATURE_NAME=demo \
  METRIC_NAMESPACE=demo_exporter \
  DEFAULT_PORT=9888 \
  TARGET_DIR=/tmp/prometheus-demo-exporter
```

Then validate the generated repository:

```bash
cd /tmp/prometheus-demo-exporter
go mod tidy
make go-check
make check
```

`FEATURE_NAME`, `METRIC_NAMESPACE`, `PROJECT_DESC`, `GO_MODULE`, and
`DEFAULT_PORT` have defaults, but passing them explicitly keeps the generated
repository predictable.

`TARGET_DIR` defaults to `rendered/$(PROJECT_NAME)` for local experiments.
Run `make check` in this scaffold repository to render a demo exporter, check
for unresolved placeholders, verify scaffold drift, and run generated Go-only
checks.

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
make drift-check TARGET_DIR=../prometheus-demo-exporter
```

To update the default managed files:

```bash
make drift-sync TARGET_DIR=../prometheus-demo-exporter
```

The default managed set is intentionally conservative: CI files, ignore files,
`cmd/main.go`, Dependabot config, `Makefile`, `Makefile.mk`, and the thin
scaffold-owned adapter in `internal/exporter/exporter.go`. It also includes the
stable feature contract, stable feature config contract, and shared feature test
helpers under `internal/<feature-name>`.
Concrete exporters keep domain logic in adjacent feature-package files, so
inspect those files separately instead of blindly syncing them:

```bash
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=Makefile
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=Dockerfile
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/exporter/exporter.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/metrics.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/snapshot_types.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/feature_config_ext.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/feature_metrics.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/collector_snapshot_test.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/collector_refresh_test.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/collector_defaults_test.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/feature_snapshot.go
make drift-check TARGET_DIR=../prometheus-demo-exporter FILE=internal/demo/feature_snapshotter.go
```

Use `ALLOW_DIRTY=1` with `make drift-sync` when you intentionally want to sync
over already modified managed files. `make drift-list-files` prints the default
managed set.

`make drift-check` also compares the target exporter's
`prometheus-exporter-framework` requirement in `go.mod` with the scaffold
version from `template/go.mod`. If the target exporter uses an older framework,
the check prints an `OUTDATED framework ...` line and exits non-zero.

Older exporters may still have scaffold-owned bootstrap files under
`internal/exporter`, such as `defaults.go`, `feature.go`, `info.go`,
`standard_metrics.go`, and their tests. Current scaffold shape replaces that
set with `internal/exporter/exporter.go`; project metadata, standard metric
names, and binary smoke metadata are supplied by the framework through
Makefile-injected linker variables.
