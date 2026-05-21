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
  --description TEXT     Defaults to the first README.md H1 or project name.
  --feature-name NAME    Defaults to FeatureName() from internal/exporter/feature.go.
  --namespace NAME       Defaults to Namespace: from tests or derived from module.
  --port PORT            Defaults to defaultListenAddress from feature.go or 9888.

File selection:
  --file PATH            Compare/sync this rendered path. Can be repeated.
  --list-files           Print the default managed file list and exit.

Default managed files:
  .dockerignore
  .github/workflows/ci.yml
  .gitignore
  .gitlab-ci.yml
  cmd/main.go

Makefiles often contain domain-specific smoke-test commands in concrete
exporters. Inspect them with --file Makefile and port relevant hunks manually.
Dockerfiles can also be domain-specific when exporters need runtime packages.
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
  ".dockerignore"
  ".github/workflows/ci.yml"
  ".gitignore"
  ".gitlab-ci.yml"
  "cmd/main.go"
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
  [[ -f "$target_dir/Makefile" ]] || return 0
  awk -F '\\?=' '
    /^[[:space:]]*PROJECT_NAME[[:space:]]*\?=/ {
      value = $2
      sub(/^[[:space:]]*/, "", value)
      sub(/[[:space:]]*$/, "", value)
      print value
      exit
    }
  ' "$target_dir/Makefile"
}

detect_readme_h1() {
  [[ -f "$target_dir/README.md" ]] || return 0
  awk '/^#[[:space:]]+/ {sub(/^#[[:space:]]+/, ""); print; exit}' "$target_dir/README.md"
}

detect_feature_name() {
  local file="$target_dir/internal/exporter/feature.go"
  [[ -f "$file" ]] || return 0
  awk '
    /FeatureName\(\)[[:space:]]+string/ {in_func = 1}
    in_func && /return[[:space:]]+"/ {
      line = $0
      sub(/^.*return[[:space:]]+"/, "", line)
      sub(/".*$/, "", line)
      print line
      exit
    }
  ' "$file"
}

detect_default_port() {
  local file="$target_dir/internal/exporter/feature.go"
  [[ -f "$file" ]] || return 0
  awk '
    /defaultListenAddress[[:space:]]*=/ && /":[0-9]+"/ {
      line = $0
      sub(/^.*:"?/, "", line)
      sub(/".*$/, "", line)
      print line
      exit
    }
  ' "$file"
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
      sed -n 's/.*Namespace:[[:space:]]*"\([^"]*\)".*/\1/p' "$file"
    done | head -n 1)"
  fi
  if [[ -n "$match" ]]; then
    printf '%s' "$match"
    return 0
  fi
  derive_namespace_from_project "${go_module:-$project_name}"
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
