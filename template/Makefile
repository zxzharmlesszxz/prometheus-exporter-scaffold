GO ?= go
GOFMT ?= gofmt
DOCKER ?= docker
STATICCHECK_VERSION ?= v0.7.0
STATICCHECK ?= $(GO) run honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)
COVERAGE_PROFILE ?= coverage.out
COVERAGE_REPORT ?= coverage.txt
COVERAGE_THRESHOLD ?= 85.0
PROJECT_NAME ?= __PROJECT_NAME__
MAIN_PACKAGE ?= ./cmd
DIST_DIR ?= dist
BUILD_OUTPUT ?= $(DIST_DIR)/$(PROJECT_NAME)
CGO_ENABLED ?= 0
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo dev)
REVISION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
BUILD_USER ?= $(shell whoami 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS ?= -s -w -X github.com/prometheus/common/version.Version=$(VERSION) -X github.com/prometheus/common/version.Branch=$(BRANCH) -X github.com/prometheus/common/version.Revision=$(REVISION) -X github.com/prometheus/common/version.BuildUser=$(BUILD_USER) -X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE)
PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
IMAGE_TAG ?= $(VERSION)
DOCKER_IMAGE ?= $(PROJECT_NAME):$(IMAGE_TAG)
DOCKER_ENTRYPOINT_NAME ?= exporter
DOCKER_HTTP_IMAGE ?= busybox:1.36
DOCKER_PLATFORMS ?= linux/amd64,linux/arm64
PROMTOOL_IMAGE ?= prom/prometheus:latest
DOCKER_COMPOSE ?= $(DOCKER) compose
SMOKE_DOCKER_IMAGE ?= $(PROJECT_NAME):smoke
SMOKE_VERSION ?= v9.8.7
SMOKE_BRANCH ?= smoke-branch
SMOKE_REVISION ?= abc123def
SMOKE_BUILD_USER ?= smoke-test
SMOKE_BUILD_DATE ?= 2026-05-17T00:00:00Z

.PHONY: help build release release-archives release-checksums release-smoke docker-build docker-buildx docker-buildx-push docker-push fmt fmt-check vet staticcheck test test-race coverage coverage-check smoke promtool-check compose-config examples-check docker-smoke-build docker-smoke-image docker-smoke go-check check full-check clean

help: ## Show available make targets.
	@printf "\033[33mUsage:\033[0m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "};{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fmt: ## Format Go files.
	$(GOFMT) -w $$($(GO) list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}{{range .TestGoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./...)

fmt-check: ## Check Go formatting.
	@test -z "$$($(GOFMT) -l $$($(GO) list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}{{range .TestGoFiles}}{{$$.Dir}}/{{.}} {{end}}' ./... | tr '\n' ' '))"

build: ## Build the local release binary into dist/.
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BUILD_OUTPUT) $(MAIN_PACKAGE)

release-archives: ## Cross-build release archives into dist/.
	mkdir -p $(DIST_DIR)
	@set -e; \
	for platform in $(PLATFORMS); do \
		goos="$${platform%/*}"; \
		goarch="$${platform#*/}"; \
		archive="$(PROJECT_NAME)_$(VERSION)_$${goos}_$${goarch}"; \
		workdir="$(DIST_DIR)/$$archive"; \
		binary="$(PROJECT_NAME)"; \
		if [ "$$goos" = "windows" ]; then binary="$$binary.exe"; fi; \
		rm -rf "$$workdir"; \
		mkdir -p "$$workdir"; \
		echo "building $$archive"; \
		CGO_ENABLED=$(CGO_ENABLED) GOOS="$$goos" GOARCH="$$goarch" $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o "$$workdir/$$binary" $(MAIN_PACKAGE); \
		cp README.md METRICS.md "$$workdir/"; \
		cp -R examples "$$workdir/examples"; \
		tar -C $(DIST_DIR) -czf "$(DIST_DIR)/$$archive.tar.gz" "$$archive"; \
		rm -rf "$$workdir"; \
	done

release-checksums: release-archives ## Write SHA256 checksums for release archives.
	@set -e; \
	cd $(DIST_DIR); \
	if command -v sha256sum >/dev/null 2>&1; then \
		sha256sum *.tar.gz > checksums.txt; \
	else \
		shasum -a 256 *.tar.gz > checksums.txt; \
	fi; \
	cat checksums.txt

release: clean release-checksums ## Build release archives and checksums.

release-smoke: release ## Build release archives and smoke-test the native archive.
	@set -e; \
	goos="$$( $(GO) env GOOS )"; \
	goarch="$$( $(GO) env GOARCH )"; \
	archive="$(DIST_DIR)/$(PROJECT_NAME)_$(VERSION)_$${goos}_$${goarch}.tar.gz"; \
	if [ ! -f "$$archive" ]; then \
		echo "skipping release smoke: native archive $$archive was not built"; \
		exit 0; \
	fi; \
	tmp="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmp"' EXIT; \
	tar -C "$$tmp" -xzf "$$archive"; \
	binary="$$tmp/$(PROJECT_NAME)_$(VERSION)_$${goos}_$${goarch}/$(PROJECT_NAME)"; \
	if [ "$$goos" = "windows" ]; then binary="$$binary.exe"; fi; \
	"$$binary" --help 2>&1 | grep -F "usage: $(PROJECT_NAME) [<flags>]" >/dev/null; \
	"$$binary" --version 2>&1 | grep -F "$(VERSION)" >/dev/null

vet: ## Run go vet.
	$(GO) vet ./...

staticcheck: ## Run staticcheck.
	$(STATICCHECK) ./...

test: ## Run Go tests.
	$(GO) test ./...

test-race: ## Run Go tests with the race detector.
	$(GO) test -race ./...

coverage: ## Run tests with coverage and write coverage reports.
	$(GO) test ./... -covermode=atomic -coverprofile=$(COVERAGE_PROFILE)
	$(GO) tool cover -func=$(COVERAGE_PROFILE) | tee $(COVERAGE_REPORT)

coverage-check: coverage ## Enforce the coverage threshold.
	@coverage="$$(awk '/^total:/ {gsub(/%/, "", $$3); print $$3}' $(COVERAGE_REPORT))"; \
	awk -v coverage="$$coverage" -v threshold="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		if (coverage + 0 < threshold + 0) { \
			printf "coverage %.1f%% is below %.1f%%\n", coverage, threshold; \
			exit 1; \
		} \
		printf "coverage %.1f%% meets threshold %.1f%%\n", coverage, threshold; \
	}'

smoke: ## Build and smoke-test the local binary.
	RUN_BINARY_SMOKE=1 GO="$(GO)" $(GO) test ./smoke -run TestBinarySmoke -count=1

docker-build: ## Build the Docker image.
	$(DOCKER) build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BRANCH=$(BRANCH) \
		--build-arg REVISION=$(REVISION) \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg PROJECT_NAME=$(PROJECT_NAME) \
		-t $(DOCKER_IMAGE) \
		.

docker-buildx: ## Build a multi-platform Docker image with buildx.
	$(DOCKER) buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BRANCH=$(BRANCH) \
		--build-arg REVISION=$(REVISION) \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg PROJECT_NAME=$(PROJECT_NAME) \
		-t $(DOCKER_IMAGE) \
		.

docker-buildx-push: ## Build and push a multi-platform Docker image with buildx.
	$(DOCKER) buildx build \
		--push \
		--platform $(DOCKER_PLATFORMS) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BRANCH=$(BRANCH) \
		--build-arg REVISION=$(REVISION) \
		--build-arg BUILD_USER=$(BUILD_USER) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg PROJECT_NAME=$(PROJECT_NAME) \
		-t $(DOCKER_IMAGE) \
		.

docker-push: ## Push the Docker image.
	$(DOCKER) push $(DOCKER_IMAGE)

promtool-check: ## Validate Prometheus example config and rules.
	$(DOCKER) run --rm --entrypoint promtool -v "$(CURDIR)/examples/prometheus:/etc/prometheus:ro" $(PROMTOOL_IMAGE) check config /etc/prometheus/prometheus.yml
	$(DOCKER) run --rm --entrypoint sh -v "$(CURDIR)/examples/prometheus:/etc/prometheus:ro" $(PROMTOOL_IMAGE) -c 'promtool check rules /etc/prometheus/rules/*.yml'

compose-config: ## Validate the Docker Compose example.
	$(DOCKER_COMPOSE) config >/dev/null

examples-check: promtool-check compose-config ## Validate shipped Prometheus and Compose examples.

docker-smoke-build: ## Build the Docker image used by docker-smoke.
	$(MAKE) docker-build \
		DOCKER_IMAGE=$(SMOKE_DOCKER_IMAGE) \
		VERSION=$(SMOKE_VERSION) \
		BRANCH=$(SMOKE_BRANCH) \
		REVISION=$(SMOKE_REVISION) \
		BUILD_USER=$(SMOKE_BUILD_USER) \
		BUILD_DATE=$(SMOKE_BUILD_DATE)

docker-smoke-image: ## Smoke-test an already built Docker image.
	@help_output="$$( $(DOCKER) run --rm $(DOCKER_IMAGE) --help 2>&1 )"; \
	echo "$$help_output" | grep -F "usage: $(DOCKER_ENTRYPOINT_NAME) [<flags>]" >/dev/null || { echo "$$help_output"; exit 1; }
	@version_output="$$( $(DOCKER) run --rm $(DOCKER_IMAGE) --version 2>&1 )"; \
	echo "$$version_output"; \
	echo "$$version_output" | grep -F "$(VERSION)" >/dev/null; \
	echo "$$version_output" | grep -F "$(BRANCH)" >/dev/null; \
	echo "$$version_output" | grep -F "$(REVISION)" >/dev/null
	@cid="$$( $(DOCKER) run -d --rm $(DOCKER_IMAGE) --log.level=error --web.listen-address=:9900 --__FEATURE_NAME__.refresh-interval=100ms )"; \
	trap "$(DOCKER) rm -f $$cid >/dev/null 2>&1 || true" EXIT; \
	i=0; \
	while [ "$$i" -lt 60 ]; do \
		i=$$((i + 1)); \
		if $(DOCKER) run --rm --network container:$$cid $(DOCKER_HTTP_IMAGE) wget -qO- http://127.0.0.1:9900/healthz 2>/dev/null | grep -qx 'ok'; then \
			break; \
		fi; \
		if [ "$$i" -eq 60 ]; then \
			$(DOCKER) logs "$$cid"; \
			exit 1; \
		fi; \
		sleep 1; \
	done; \
	metrics=""; \
	i=0; \
	while [ "$$i" -lt 60 ]; do \
		i=$$((i + 1)); \
		metrics="$$( $(DOCKER) run --rm --network container:$$cid $(DOCKER_HTTP_IMAGE) wget -qO- http://127.0.0.1:9900/metrics )"; \
		if echo "$$metrics" | grep -F '__METRIC_NAMESPACE___last_collection_success 1' >/dev/null; then \
			break; \
		fi; \
		if [ "$$i" -eq 60 ]; then \
			$(DOCKER) logs "$$cid"; \
			echo "$$metrics"; \
			exit 1; \
		fi; \
		sleep 1; \
	done; \
	echo "$$metrics" | grep -F "__METRIC_NAMESPACE___build_info" >/dev/null; \
	echo "$$metrics" | grep -F 'version="$(VERSION)"' >/dev/null; \
	echo "$$metrics" | grep -F 'branch="$(BRANCH)"' >/dev/null; \
	echo "$$metrics" | grep -F 'revision="$(REVISION)"' >/dev/null; \
	echo "$$metrics" | grep -F '__FEATURE_NAME___example_value 1' >/dev/null; \
	$(DOCKER) rm -f "$$cid" >/dev/null 2>&1 || true

docker-smoke: docker-smoke-build ## Build and smoke-test the Docker image.
	$(MAKE) docker-smoke-image \
		DOCKER_IMAGE=$(SMOKE_DOCKER_IMAGE) \
		VERSION=$(SMOKE_VERSION) \
		BRANCH=$(SMOKE_BRANCH) \
		REVISION=$(SMOKE_REVISION)

go-check: fmt-check vet staticcheck coverage-check smoke test-race ## Run Go checks that do not require Docker.

check: go-check examples-check ## Run the standard maintenance check.

full-check: check docker-smoke release-smoke ## Run all local checks, Docker smoke, and release smoke.

clean: ## Remove generated local artifacts.
	rm -rf $(DIST_DIR)
	rm -f $(COVERAGE_PROFILE) $(COVERAGE_REPORT)
