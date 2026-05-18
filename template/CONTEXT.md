# Context For Next Chat

This repository is `__PROJECT_NAME__`.

It was created from:

`prometheus-exporter-git-template`

## Template Exporter Context

`prometheus-template-exporter` provides the reusable exporter shell:

- CLI bootstrap through `exporter.MainForProject(...)`
- standard Prometheus/exporter-toolkit flags
- `promslog` logging
- `/metrics`, `/healthz`, landing page, optional pprof
- custom Prometheus registry
- build info, Go runtime, and process collectors
- version metadata hydration from linker flags or Go build info

The template module is consumed from its published module path:

```go
github.com/zxzharmlesszxz/prometheus-template-exporter v0.1.2
```

## Current Exporter State

Entrypoint:

`cmd/main.go`

```go
template.MainForProject(
	"__PROJECT_NAME__",
	"__PROJECT_DESC__",
	exporter.NewFeature(),
)
```

Default listen address:

`:__DEFAULT_PORT__`

Default refresh interval:

`1m`

## Verification Targets

```bash
make help
make go-check
make check
make docker-smoke
make full-check
```
