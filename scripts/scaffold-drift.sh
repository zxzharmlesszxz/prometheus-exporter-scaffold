#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter
  scripts/scaffold-drift.sh --target-dir ../prometheus-demo-exporter --sync

Checks or syncs scaffold-owned files in an existing exporter against the current
template. By default the script only reports drift and exits non-zero when drift
is found.

Required:
  --target-dir DIR       Existing exporter repository to compare.

Modes:
  --check                Report drift only. This is the default.
  --sync                 Copy rendered scaffold-owned files into --target-dir.
  --allow-dirty          Allow --sync when managed files already have git changes.

Render metadata overrides:
  --project-name NAME    Defaults to basename of --target-dir.
  --module MODULE        Defaults to module path from go.mod.
  --description TEXT     Defaults to rendered exporter description, README H1, or project name.
  --feature-name NAME    Defaults to DefaultFeatureName/defaultFeatureName or FeatureName().
  --namespace NAME       Defaults to DefaultMetricNamespace, Namespace: from tests, or derived.
  --port PORT            Defaults to DefaultListenAddress/defaultListenAddress or 9888.

File selection:
  --file PATH            Compare/sync this rendered path. Can be repeated.
  --list-files           Print the default managed file list and exit.

Default managed files:
  exporter.mk
  .dockerignore
  .github/workflows/ci.yml
  .gitignore
  .gitlab-ci.yml
  cmd/main.go
  internal/exporter/defaults.go
  internal/exporter/identity.go
  internal/exporter/main.go
  internal/exporter/standard_metrics.go
  internal/exporter/variables/variables.go

Makefiles often contain domain-specific smoke-test commands in concrete
exporters. Inspect target logic with --file Makefile and port relevant hunks
manually. Common make variables live in exporter.mk and are scaffold-managed.
Dockerfiles can also be domain-specific when exporters need runtime packages.
Legacy exporters may still define Main(), FeatureName(), or
DefaultListenAddress() in internal/exporter/feature.go, or keep rendered
defaults in older files. Remove those definitions once when adopting the split
scaffold Go files.
Metric constants are split so scaffold-owned standard names live in
internal/exporter/standard_metrics.go. Domain-specific metric constants should
remain in internal/exporter/metrics.go.
Collector types are domain-specific; inspect split collector type drift with
--file internal/exporter/collector_types.go before migrating manually.
Collector metric descriptors and emission are domain-specific; inspect split
collector metric method drift with --file internal/exporter/collector_metrics.go.
Snapshot gatherers and snapshot status/error adapters are also domain-specific;
inspect them with --file internal/exporter/snapshot.go.
Feature flags, collector registration, and runtime config are domain-specific;
inspect split feature method drift with --file internal/exporter/feature_flags.go,
--file internal/exporter/feature_collectors.go, or
--file internal/exporter/runtime_config.go.
Collector test helpers can be inspected with
--file internal/exporter/collector_test_helpers_test.go.
Collector tests can be inspected by concern with
--file internal/exporter/collector_snapshot_test.go,
--file internal/exporter/collector_refresh_test.go, or
--file internal/exporter/collector_defaults_test.go.
Feature test helpers can be inspected with
--file internal/exporter/feature_test_helpers_test.go and
--file internal/exporter/feature_integration_test_helpers_test.go.
Feature tests can be inspected by concern with
--file internal/exporter/feature_flags_test.go,
--file internal/exporter/feature_collectors_test.go,
--file internal/exporter/runtime_config_test.go, or
--file internal/exporter/identity_test.go.
USAGE
}

mode="check"
allow_dirty=0
list_files=0
target_dir=""
project_name=""
go_module=""
project_desc=""
feature_name=""
metric_namespace=""
default_port=""
custom_files=()

default_files=(
  "exporter.mk"
  ".dockerignore"
  ".github/workflows/ci.yml"
  ".gitignore"
  ".gitlab-ci.yml"
  "cmd/main.go"
  "internal/exporter/defaults.go"
  "internal/exporter/identity.go"
  "internal/exporter/main.go"
  "internal/exporter/standard_metrics.go"
  "internal/exporter/variables/variables.go"
)

while [[ $# -gt 0 ]]; do
  case "$1" in
    --target-dir|--exporter-dir)
      target_dir="${2:-}"
      shift 2
      ;;
    --check)
      mode="check"
      shift
      ;;
    --sync)
      mode="sync"
      shift
      ;;
    --allow-dirty)
      allow_dirty=1
      shift
      ;;
    --project-name)
      project_name="${2:-}"
      shift 2
      ;;
    --module)
      go_module="${2:-}"
      shift 2
      ;;
    --description)
      project_desc="${2:-}"
      shift 2
      ;;
    --feature-name)
      feature_name="${2:-}"
      shift 2
      ;;
    --namespace)
      metric_namespace="${2:-}"
      shift 2
      ;;
    --port)
      default_port="${2:-}"
      shift 2
      ;;
    --file)
      custom_files+=("${2:-}")
      shift 2
      ;;
    --list-files)
      list_files=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ "$list_files" -eq 1 ]]; then
  printf '%s\n' "${default_files[@]}"
  exit 0
fi

if [[ -z "$target_dir" ]]; then
  usage >&2
  exit 1
fi
if [[ ! -d "$target_dir" ]]; then
  echo "target dir does not exist: $target_dir" >&2
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd "$script_dir/.." && pwd)"
target_dir="$(cd "$target_dir" && pwd)"

managed_files=("${default_files[@]}")
if [[ "${#custom_files[@]}" -gt 0 ]]; then
  managed_files=("${custom_files[@]}")
fi

detect_module() {
  [[ -f "$target_dir/go.mod" ]] || return 0
  awk '$1 == "module" {print $2; exit}' "$target_dir/go.mod"
}

detect_project_name() {
  local file value
  for file in "$target_dir/exporter.mk" "$target_dir/Makefile"; do
    [[ -f "$file" ]] || continue
    value="$(awk -F '\\?=' '
      /^[[:space:]]*PROJECT_NAME[[:space:]]*\?=/ {
        value = $2
        sub(/^[[:space:]]*/, "", value)
        sub(/[[:space:]]*$/, "", value)
        print value
        exit
      }
    ' "$file")"
    if [[ -n "$value" && "$value" != *'$('* ]]; then
      printf '%s' "$value"
      return 0
    fi
  done
  if [[ -d "$target_dir/internal/exporter" ]]; then
    while IFS= read -r file; do
      value="$(awk '
        /DefaultExporterName[[:space:]]*=/ && /"/ {
          line = $0
          sub(/^.*DefaultExporterName[[:space:]]*=[[:space:]]*"/, "", line)
          sub(/".*$/, "", line)
          print line
          exit
        }
      ' "$file")"
      if [[ -n "$value" ]]; then
        printf '%s' "$value"
        return 0
      fi
    done < <(find "$target_dir/internal/exporter" -type f -name '*.go' -print 2>/dev/null | sort)
  fi
}

detect_readme_h1() {
  [[ -f "$target_dir/README.md" ]] || return 0
  awk '/^#[[:space:]]+/ {sub(/^#[[:space:]]+/, ""); print; exit}' "$target_dir/README.md"
}

detect_exporter_description() {
  local dir="$target_dir/internal/exporter"
  local file value
  [[ -d "$dir" ]] || return 0
  while IFS= read -r file; do
    value="$(awk '
      /DefaultExporterDescription[[:space:]]*=/ && /"/ {
        line = $0
        sub(/^.*DefaultExporterDescription[[:space:]]*=[[:space:]]*"/, "", line)
        sub(/".*$/, "", line)
        print line
        exit
      }
    ' "$file")"
    if [[ -n "$value" ]]; then
      printf '%s' "$value"
      return 0
    fi
  done < <(find "$dir" -type f -name '*.go' -print 2>/dev/null | sort)
}

detect_feature_name() {
  local dir="$target_dir/internal/exporter"
  local file value
  [[ -d "$dir" ]] || return 0
  while IFS= read -r file; do
    value="$(awk '
      /defaultFeatureName[[:space:]]*=/ && /"/ {
        line = $0
        sub(/^.*defaultFeatureName[[:space:]]*=[[:space:]]*"/, "", line)
        sub(/".*$/, "", line)
        print line
        exit
      }
      /DefaultFeatureName[[:space:]]*=/ && /"/ {
        line = $0
        sub(/^.*DefaultFeatureName[[:space:]]*=[[:space:]]*"/, "", line)
        sub(/".*$/, "", line)
        print line
        exit
      }
      /FeatureName\(\)[[:space:]]+string/ {in_func = 1}
      in_func && /return[[:space:]]+"/ {
        line = $0
        sub(/^.*return[[:space:]]+"/, "", line)
        sub(/".*$/, "", line)
        print line
        exit
      }
      in_func && /^}/ {in_func = 0}
    ' "$file")"
    if [[ -n "$value" ]]; then
      printf '%s' "$value"
      return 0
    fi
  done < <(find "$dir" -type f -name '*.go' -print 2>/dev/null | sort)
}

detect_default_port() {
  local dir="$target_dir/internal/exporter"
  local file value
  [[ -d "$dir" ]] || return 0
  while IFS= read -r file; do
    value="$(awk '
      /defaultListenAddress[[:space:]]*=/ && /":[0-9]+"/ {
        line = $0
        sub(/^.*:"?/, "", line)
        sub(/".*$/, "", line)
        print line
        exit
      }
      /DefaultListenAddress[[:space:]]*=/ && /":[0-9]+"/ {
        line = $0
        sub(/^.*:"?/, "", line)
        sub(/".*$/, "", line)
        print line
        exit
      }
    ' "$file")"
    if [[ -n "$value" ]]; then
      printf '%s' "$value"
      return 0
    fi
  done < <(find "$dir" -type f -name '*.go' -print 2>/dev/null | sort)
}

sanitize_metric_namespace() {
  local value="$1"
  value="$(printf '%s' "$value" | tr '[:upper:]' '[:lower:]' | sed -e 's/[^a-z0-9]/_/g' -e 's/_\{1,\}/_/g' -e 's/^_*//' -e 's/_*$//')"
  if [[ -z "$value" ]]; then
    value="exporter_framework"
  fi
  if [[ "$value" =~ ^[0-9] ]]; then
    value="_$value"
  fi
  if [[ "$value" != "exporter_framework" && "$value" != *_exporter ]]; then
    value="${value}_exporter"
  fi
  printf '%s' "$value"
}

derive_namespace_from_project() {
  local project="$1"
  local base="${project##*/}"
  base="${base#prometheus-}"
  base="${base%-exporter}"
  sanitize_metric_namespace "$base"
}

detect_namespace() {
  local match=""
  if [[ -d "$target_dir/internal/exporter" ]]; then
    match="$(find "$target_dir/internal/exporter" -type f -name '*.go' -print 2>/dev/null | sort | while IFS= read -r file; do
      sed -n \
        -e 's/.*DefaultMetricNamespace[[:space:]]*=[[:space:]]*"\([^"]*\)".*/\1/p' \
        -e 's/.*defaultMetricNamespace[[:space:]]*=[[:space:]]*"\([^"]*\)".*/\1/p' \
        -e 's/.*Namespace:[[:space:]]*"\([^"]*\)".*/\1/p' \
        "$file"
    done | head -n 1)"
  fi
  if [[ -n "$match" ]]; then
    printf '%s' "$match"
    return 0
  fi
  derive_namespace_from_project "${go_module:-$project_name}"
}

feature_go_defines() {
  local pattern="$1"
  local file="$target_dir/internal/exporter/feature.go"
  [[ -f "$file" ]] || return 1
  grep -Eq "$pattern" "$file"
}

exporter_go_defines_except() {
  local skip_path="$target_dir/$1"
  local pattern="$2"
  local dir="$target_dir/internal/exporter"
  local file
  [[ -d "$dir" ]] || return 1
  while IFS= read -r file; do
    if [[ "$file" == "$skip_path" ]]; then
      continue
    fi
    if grep -Eq "$pattern" "$file"; then
      return 0
    fi
  done < <(find "$dir" -type f -name '*.go' -print 2>/dev/null | sort)
  return 1
}

legacy_managed_go_reason() {
  local file="$1"
  case "$file" in
    internal/exporter/main.go)
      if feature_go_defines '^[[:space:]]*func[[:space:]]+Main[[:space:]]*\('; then
        echo "Main() is still defined in internal/exporter/feature.go"
        return 0
      fi
      ;;
    internal/exporter/identity.go)
      local reasons=()
      if feature_go_defines '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+FeatureName[[:space:]]*\('; then
        reasons+=("FeatureName()")
      fi
      if feature_go_defines '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+DefaultListenAddress[[:space:]]*\('; then
        reasons+=("DefaultListenAddress()")
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/feature.go"
        return 0
      fi
      ;;
    internal/exporter/defaults.go)
      local reasons=()
      if exporter_go_defines_except "$file" '^[[:space:]]*defaultExporterName[[:space:]]*='; then
        reasons+=("defaultExporterName")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*defaultExporterDescription[[:space:]]*='; then
        reasons+=("defaultExporterDescription")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*defaultFeatureName[[:space:]]*='; then
        reasons+=("defaultFeatureName")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*defaultMetricNamespace[[:space:]]*='; then
        reasons+=("defaultMetricNamespace")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*defaultListenAddress[[:space:]]*='; then
        reasons+=("defaultListenAddress")
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined outside internal/exporter/defaults.go"
        return 0
      fi
      ;;
    internal/exporter/standard_metrics.go)
      local reasons=()
      if exporter_go_defines_except "$file" '^[[:space:]]*metricBuildInfo[[:space:]]*='; then
        reasons+=("metricBuildInfo")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*metricLastCollectionSuccess[[:space:]]*='; then
        reasons+=("metricLastCollectionSuccess")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*metricLastCollectionTimestampSeconds[[:space:]]*='; then
        reasons+=("metricLastCollectionTimestampSeconds")
      fi
      if exporter_go_defines_except "$file" '^[[:space:]]*metricLastSuccessfulCollectionTimestampSeconds[[:space:]]*='; then
        reasons+=("metricLastSuccessfulCollectionTimestampSeconds")
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined outside internal/exporter/standard_metrics.go"
        return 0
      fi
      ;;
    internal/exporter/feature_flags.go)
      if feature_go_defines '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+RegisterFlags[[:space:]]*\('; then
        echo "RegisterFlags() still defined in internal/exporter/feature.go"
        return 0
      fi
      ;;
    internal/exporter/feature_collectors.go)
      if feature_go_defines '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+RegisterCollectors[[:space:]]*\('; then
        echo "RegisterCollectors() still defined in internal/exporter/feature.go"
        return 0
      fi
      ;;
    internal/exporter/runtime_config.go)
      if feature_go_defines '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+RuntimeConfig[[:space:]]*\('; then
        echo "RuntimeConfig() still defined in internal/exporter/feature.go"
        return 0
      fi
      ;;
    internal/exporter/collector_types.go)
      local reasons=()
      local collector_file="$target_dir/internal/exporter/collector.go"
      if [[ -f "$collector_file" ]]; then
        if grep -Eq '^[[:space:]]*type[[:space:]]+Snapshot([[:space:]]|$)' "$collector_file"; then
          reasons+=("Snapshot")
        fi
        if grep -Eq '^[[:space:]]*type[[:space:]]+SnapshotGatherer([[:space:]]|$)' "$collector_file"; then
          reasons+=("SnapshotGatherer")
        fi
        if grep -Eq '^[[:space:]]*type[[:space:]]+Collector([[:space:]]|$)' "$collector_file"; then
          reasons+=("Collector")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/collector.go"
        return 0
      fi
      ;;
    internal/exporter/collector_metrics.go)
      local reasons=()
      local collector_file="$target_dir/internal/exporter/collector.go"
      if [[ -f "$collector_file" ]]; then
        if grep -Eq '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+describeSnapshotMetrics[[:space:]]*\(' "$collector_file"; then
          reasons+=("describeSnapshotMetrics()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+collectSnapshotMetrics[[:space:]]*\(' "$collector_file"; then
          reasons+=("collectSnapshotMetrics()")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/collector.go"
        return 0
      fi
      ;;
    internal/exporter/collector_test_helpers_test.go)
      local reasons=()
      local collector_test_file="$target_dir/internal/exporter/collector_test.go"
      if [[ -f "$collector_test_file" ]]; then
        if grep -Eq '^[[:space:]]*type[[:space:]]+fakeSnapshotter([[:space:]]|$)' "$collector_test_file"; then
          reasons+=("fakeSnapshotter")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]+newFakeSnapshotter[[:space:]]*\(' "$collector_test_file"; then
          reasons+=("newFakeSnapshotter()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+Snapshot[[:space:]]*\(' "$collector_test_file"; then
          reasons+=("fakeSnapshotter.Snapshot()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]*\([^)]*\)[[:space:]]+set[[:space:]]*\(' "$collector_test_file"; then
          reasons+=("fakeSnapshotter.set()")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/collector_test.go"
        return 0
      fi
      ;;
    internal/exporter/collector_snapshot_test.go)
      local collector_test_file="$target_dir/internal/exporter/collector_test.go"
      if [[ -f "$collector_test_file" ]] && grep -Eq '^[[:space:]]*func[[:space:]]+TestCollectorExportsSnapshot[[:space:]]*\(' "$collector_test_file"; then
        echo "TestCollectorExportsSnapshot() still defined in internal/exporter/collector_test.go"
        return 0
      fi
      ;;
    internal/exporter/collector_refresh_test.go)
      local collector_test_file="$target_dir/internal/exporter/collector_test.go"
      if [[ -f "$collector_test_file" ]] && grep -Eq '^[[:space:]]*func[[:space:]]+TestCollectorBackgroundRefresh[^[:space:]]*[[:space:]]*\(' "$collector_test_file"; then
        echo "TestCollectorBackgroundRefresh*() still defined in internal/exporter/collector_test.go"
        return 0
      fi
      ;;
    internal/exporter/collector_defaults_test.go)
      local reasons=()
      local collector_test_file="$target_dir/internal/exporter/collector_test.go"
      if [[ -f "$collector_test_file" ]]; then
        if grep -Eq '^[[:space:]]*func[[:space:]]+TestCollectorDefaults[^[:space:]]*[[:space:]]*\(' "$collector_test_file"; then
          reasons+=("TestCollectorDefaults*()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]+TestCollectorUsesDefaultSnapshotter[[:space:]]*\(' "$collector_test_file"; then
          reasons+=("TestCollectorUsesDefaultSnapshotter()")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/collector_test.go"
        return 0
      fi
      ;;
    internal/exporter/feature_test_helpers_test.go)
      local feature_test_file="$target_dir/internal/exporter/feature_test.go"
      if [[ -f "$feature_test_file" ]] && grep -Eq '^[[:space:]]*func[[:space:]]+testFeatureContext[[:space:]]*\(' "$feature_test_file"; then
        echo "testFeatureContext() still defined in internal/exporter/feature_test.go"
        return 0
      fi
      ;;
    internal/exporter/feature_integration_test_helpers_test.go)
      local reasons=()
      local feature_integration_test_file="$target_dir/internal/exporter/feature_integration_test.go"
      if [[ -f "$feature_integration_test_file" ]]; then
        if grep -Eq '^[[:space:]]*func[[:space:]]+newTestHandler[[:space:]]*\(' "$feature_integration_test_file"; then
          reasons+=("newTestHandler()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]+waitForHandlerMetrics[[:space:]]*\(' "$feature_integration_test_file"; then
          reasons+=("waitForHandlerMetrics()")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/feature_integration_test.go"
        return 0
      fi
      ;;
    internal/exporter/feature_flags_test.go)
      local feature_test_file="$target_dir/internal/exporter/feature_test.go"
      if [[ -f "$feature_test_file" ]] && grep -Eq '^[[:space:]]*func[[:space:]]+TestFeatureRegistersAndParsesFlags[[:space:]]*\(' "$feature_test_file"; then
        echo "TestFeatureRegistersAndParsesFlags() still defined in internal/exporter/feature_test.go"
        return 0
      fi
      ;;
    internal/exporter/feature_collectors_test.go)
      local reasons=()
      local feature_test_file="$target_dir/internal/exporter/feature_test.go"
      if [[ -f "$feature_test_file" ]]; then
        if grep -Eq '^[[:space:]]*func[[:space:]]+TestFeatureRegistersCollector[[:space:]]*\(' "$feature_test_file"; then
          reasons+=("TestFeatureRegistersCollector()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]+TestFeatureReportsCollectorRegistrationError[[:space:]]*\(' "$feature_test_file"; then
          reasons+=("TestFeatureReportsCollectorRegistrationError()")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/feature_test.go"
        return 0
      fi
      ;;
    internal/exporter/runtime_config_test.go)
      local feature_test_file="$target_dir/internal/exporter/feature_test.go"
      if [[ -f "$feature_test_file" ]] && grep -Eq '^[[:space:]]*func[[:space:]]+TestFeatureRuntimeConfig[^[:space:]]*[[:space:]]*\(' "$feature_test_file"; then
        echo "TestFeatureRuntimeConfig*() still defined in internal/exporter/feature_test.go"
        return 0
      fi
      ;;
    internal/exporter/identity_test.go)
      local feature_test_file="$target_dir/internal/exporter/feature_test.go"
      if [[ -f "$feature_test_file" ]] && grep -Eq '^[[:space:]]*func[[:space:]]+TestFeatureMetadata[[:space:]]*\(' "$feature_test_file"; then
        echo "TestFeatureMetadata() still defined in internal/exporter/feature_test.go"
        return 0
      fi
      ;;
    internal/exporter/snapshot.go)
      local reasons=()
      local collector_file="$target_dir/internal/exporter/collector.go"
      if [[ -f "$collector_file" ]]; then
        if grep -Eq '^[[:space:]]*func[[:space:]]*\(SnapshotGatherer\)[[:space:]]+Snapshot[[:space:]]*\(' "$collector_file"; then
          reasons+=("SnapshotGatherer.Snapshot()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]+snapshotStatus[[:space:]]*\(' "$collector_file"; then
          reasons+=("snapshotStatus()")
        fi
        if grep -Eq '^[[:space:]]*func[[:space:]]+logSnapshotError[[:space:]]*\(' "$collector_file"; then
          reasons+=("logSnapshotError()")
        fi
      fi
      if [[ "${#reasons[@]}" -gt 0 ]]; then
        echo "${reasons[*]} still defined in internal/exporter/collector.go"
        return 0
      fi
      ;;
  esac
  return 1
}

if [[ -z "$go_module" ]]; then
  go_module="$(detect_module)"
fi
if [[ -z "$project_name" ]]; then
  project_name="$(detect_project_name)"
fi
if [[ -z "$project_name" ]]; then
  project_name="$(basename "$target_dir")"
fi
if [[ -z "$project_desc" ]]; then
  project_desc="$(detect_exporter_description)"
fi
if [[ -z "$project_desc" ]]; then
  project_desc="$(detect_readme_h1)"
fi
if [[ -z "$project_desc" ]]; then
  project_desc="$project_name"
fi
if [[ -z "$feature_name" ]]; then
  feature_name="$(detect_feature_name)"
fi
if [[ -z "$feature_name" ]]; then
  stem="${project_name#prometheus-}"
  stem="${stem%-exporter}"
  feature_name="${stem//-/_}"
fi
if [[ -z "$metric_namespace" ]]; then
  metric_namespace="$(detect_namespace)"
fi
if [[ -z "$default_port" ]]; then
  default_port="$(detect_default_port)"
fi
if [[ -z "$default_port" ]]; then
  default_port="9888"
fi

rendered_dir="$(mktemp -d)"
trap 'rm -rf "$rendered_dir"' EXIT

"$repo_dir/scripts/render.sh" \
  --project-name "$project_name" \
  --module "${go_module:-$project_name}" \
  --description "$project_desc" \
  --feature-name "$feature_name" \
  --namespace "$metric_namespace" \
  --port "$default_port" \
  --target-dir "$rendered_dir" >/dev/null

printf 'Scaffold metadata:\n'
printf '  target:       %s\n' "$target_dir"
printf '  project-name: %s\n' "$project_name"
printf '  module:       %s\n' "${go_module:-$project_name}"
printf '  description:  %s\n' "$project_desc"
printf '  feature-name: %s\n' "$feature_name"
printf '  namespace:    %s\n' "$metric_namespace"
printf '  port:         %s\n' "$default_port"

if [[ "$mode" == "sync" && "$allow_dirty" -ne 1 ]] && git -C "$target_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  dirty="$(git -C "$target_dir" status --short -- "${managed_files[@]}")"
  if [[ -n "$dirty" ]]; then
    echo
    echo "managed files already have git changes; commit/stash them or pass --allow-dirty:" >&2
    echo "$dirty" >&2
    exit 1
  fi
fi

drift=0
echo
for file in "${managed_files[@]}"; do
  rendered_file="$rendered_dir/$file"
  target_file="$target_dir/$file"

  if [[ ! -e "$rendered_file" ]]; then
    echo "SKIP    $file (not rendered)"
    continue
  fi

  legacy_reason="$(legacy_managed_go_reason "$file" || true)"
  if [[ -n "$legacy_reason" ]]; then
    echo "LEGACY $file ($legacy_reason; migrate before syncing this file)"
    drift=1
    continue
  fi

  if [[ "$mode" == "sync" ]]; then
    if [[ -e "$target_file" ]] && cmp -s "$rendered_file" "$target_file"; then
      echo "OK      $file"
      continue
    fi
    mkdir -p "$(dirname "$target_file")"
    cp "$rendered_file" "$target_file"
    echo "SYNCED  $file"
    continue
  fi

  if [[ ! -e "$target_file" ]]; then
    echo "MISSING $file"
    drift=1
    continue
  fi
  if cmp -s "$rendered_file" "$target_file"; then
    echo "OK      $file"
    continue
  fi

  echo "DRIFT   $file"
  diff -u "$target_file" "$rendered_file" || true
  drift=1
done

if [[ "$mode" == "check" && "$drift" -ne 0 ]]; then
  echo
  echo "scaffold drift found; rerun with --sync to update managed files"
  exit 1
fi
