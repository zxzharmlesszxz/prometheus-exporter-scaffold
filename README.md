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

The generated `cmd/main.go` is intentionally stable. Project metadata is derived
by `prometheus-exporter-framework` from the Go module path, while the concrete
feature package owns `exporter.Main()`.

## Framework Version

`framework.version` and `template/go.mod` track the
`prometheus-exporter-framework` version used by newly generated exporters.

When a new framework tag is published, the framework release workflow opens a
pull request here to update the scaffold and verify a rendered exporter.

This repository's own CI also renders a demo exporter and runs its Go-only
checks, so scaffold pull requests validate the generated code path directly.
