GO ?= go
GOFMT ?= gofmt

PROJECT_NAME ?=
GO_MODULE ?=
PROJECT_DESC ?=
FEATURE_NAME ?=
METRIC_NAMESPACE ?=
DEFAULT_PORT ?=
TARGET_DIR ?= $(if $(PROJECT_NAME),rendered/$(PROJECT_NAME),)
FILE ?=
ALLOW_DIRTY ?= 0

CHECK_PROJECT_NAME ?= prometheus-demo-exporter
CHECK_GO_MODULE ?= github.com/zxzharmlesszxz/prometheus-demo-exporter
CHECK_PROJECT_DESC ?= Prometheus Demo Exporter
CHECK_FEATURE_NAME ?= demo
CHECK_METRIC_NAMESPACE ?= demo_exporter
CHECK_DEFAULT_PORT ?= 9888

DRIFT_ALLOW_DIRTY := $(if $(filter 1 true yes,$(ALLOW_DIRTY)),--allow-dirty,)
DRIFT_FILE_ARGS := $(foreach file,$(FILE),--file "$(file)")

.PHONY: help scripts-check tools-check render-check check new-exporter drift-check drift-sync drift-list-files clean require-project-name require-target-dir

help: ## Show available make targets.
	@printf "\033[33mUsage:\033[0m\n"
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "};{printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

scripts-check: ## Check scaffold shell scripts for syntax errors.
	@bash -n scripts/render.sh
	@bash -n scripts/scaffold-drift.sh

tools-check: ## Check local tools required by generated Go checks.
	@command -v "$(GO)" >/dev/null 2>&1 || { echo "GO command not found: $(GO)" >&2; exit 127; }
	@command -v "$(GOFMT)" >/dev/null 2>&1 || { echo "GOFMT command not found: $(GOFMT)" >&2; exit 127; }

render-check: scripts-check tools-check ## Render a demo exporter and run its Go-only checks.
	@set -e; \
	tmp="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmp"' EXIT; \
	scripts/render.sh \
		--project-name "$(CHECK_PROJECT_NAME)" \
		--module "$(CHECK_GO_MODULE)" \
		--description "$(CHECK_PROJECT_DESC)" \
		--feature-name "$(CHECK_FEATURE_NAME)" \
		--namespace "$(CHECK_METRIC_NAMESPACE)" \
		--port "$(CHECK_DEFAULT_PORT)" \
		--target-dir "$$tmp"; \
	if grep -R -n -E '__[A-Z0-9_]+__' "$$tmp"; then \
		echo "rendered exporter still contains scaffold placeholders" >&2; \
		exit 1; \
	fi; \
	scripts/scaffold-drift.sh --target-dir "$$tmp"; \
	cd "$$tmp"; \
	$(GO) mod tidy; \
	$(MAKE) go-check GO="$(GO)" GOFMT="$(GOFMT)"

check: render-check ## Run scaffold checks used by CI.

new-exporter: require-project-name require-target-dir ## Render a new exporter. Set PROJECT_NAME and optionally TARGET_DIR/GO_MODULE/PROJECT_DESC/FEATURE_NAME/METRIC_NAMESPACE/DEFAULT_PORT.
	@scripts/render.sh \
		--project-name "$(PROJECT_NAME)" \
		$(if $(GO_MODULE),--module "$(GO_MODULE)",) \
		$(if $(PROJECT_DESC),--description "$(PROJECT_DESC)",) \
		$(if $(FEATURE_NAME),--feature-name "$(FEATURE_NAME)",) \
		$(if $(METRIC_NAMESPACE),--namespace "$(METRIC_NAMESPACE)",) \
		$(if $(DEFAULT_PORT),--port "$(DEFAULT_PORT)",) \
		--target-dir "$(TARGET_DIR)"

drift-check: require-target-dir ## Check scaffold-managed files in an existing exporter. Set TARGET_DIR and optionally FILE.
	@scripts/scaffold-drift.sh \
		--target-dir "$(TARGET_DIR)" \
		$(if $(PROJECT_NAME),--project-name "$(PROJECT_NAME)",) \
		$(if $(GO_MODULE),--module "$(GO_MODULE)",) \
		$(if $(PROJECT_DESC),--description "$(PROJECT_DESC)",) \
		$(if $(FEATURE_NAME),--feature-name "$(FEATURE_NAME)",) \
		$(if $(METRIC_NAMESPACE),--namespace "$(METRIC_NAMESPACE)",) \
		$(if $(DEFAULT_PORT),--port "$(DEFAULT_PORT)",) \
		$(DRIFT_FILE_ARGS)

drift-sync: require-target-dir ## Sync scaffold-managed files into an existing exporter. Set TARGET_DIR; use ALLOW_DIRTY=1 to permit dirty managed files.
	@scripts/scaffold-drift.sh \
		--target-dir "$(TARGET_DIR)" \
		--sync \
		$(DRIFT_ALLOW_DIRTY) \
		$(if $(PROJECT_NAME),--project-name "$(PROJECT_NAME)",) \
		$(if $(GO_MODULE),--module "$(GO_MODULE)",) \
		$(if $(PROJECT_DESC),--description "$(PROJECT_DESC)",) \
		$(if $(FEATURE_NAME),--feature-name "$(FEATURE_NAME)",) \
		$(if $(METRIC_NAMESPACE),--namespace "$(METRIC_NAMESPACE)",) \
		$(if $(DEFAULT_PORT),--port "$(DEFAULT_PORT)",) \
		$(DRIFT_FILE_ARGS)

drift-list-files: ## List the default scaffold-managed files used by drift-check/drift-sync.
	@scripts/scaffold-drift.sh --list-files

clean: ## Remove local rendered exporters and temporary artifacts.
	@rm -rf rendered tmp

require-project-name:
	@test -n "$(PROJECT_NAME)" || { echo "PROJECT_NAME is required" >&2; exit 2; }

require-target-dir:
	@test -n "$(TARGET_DIR)" || { echo "TARGET_DIR is required" >&2; exit 2; }
