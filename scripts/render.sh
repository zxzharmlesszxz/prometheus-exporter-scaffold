#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  scripts/render.sh \
    --project-name prometheus-demo-exporter \
    --module prometheus-demo-exporter \
    --description "Prometheus Demo Exporter" \
    --feature-name demo \
    --namespace demo_exporter \
    --port 9888 \
    --target-dir /tmp/prometheus-demo-exporter

Required:
  --project-name
  --target-dir

Optional:
  --module       Defaults to --project-name.
  --description Defaults to --project-name.
  --feature-name Defaults to project name without prometheus- prefix and -exporter suffix, with '-' replaced by '_'.
  --namespace   Defaults to <feature-name>_exporter.
  --port        Defaults to 9888.
USAGE
}

project_name=""
go_module=""
project_desc=""
feature_name=""
metric_namespace=""
default_port="9888"
target_dir=""

while [[ $# -gt 0 ]]; do
  case "$1" in
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
    --target-dir)
      target_dir="${2:-}"
      shift 2
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

if [[ -z "$project_name" || -z "$target_dir" ]]; then
  usage >&2
  exit 1
fi

go_module="${go_module:-$project_name}"
project_desc="${project_desc:-$project_name}"
if [[ -z "$feature_name" ]]; then
  stem="${project_name#prometheus-}"
  stem="${stem%-exporter}"
  feature_name="${stem//-/_}"
fi
metric_namespace="${metric_namespace:-${feature_name}_exporter}"

if [[ ! "$feature_name" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
  echo "--feature-name must be a valid Go/Prometheus identifier fragment: $feature_name" >&2
  exit 1
fi
if [[ ! "$metric_namespace" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
  echo "--namespace must be a valid Prometheus metric namespace: $metric_namespace" >&2
  exit 1
fi
if [[ ! "$default_port" =~ ^[0-9]+$ ]]; then
  echo "--port must be numeric: $default_port" >&2
  exit 1
fi

if [[ -e "$target_dir" && -n "$(find "$target_dir" -mindepth 1 -maxdepth 1 -print -quit)" ]]; then
  echo "target dir exists and is not empty: $target_dir" >&2
  exit 1
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd "$script_dir/.." && pwd)"
template_dir="$repo_dir/template"

mkdir -p "$target_dir"
cp -R "$template_dir/." "$target_dir/"

export PROJECT_NAME="$project_name"
export GO_MODULE="$go_module"
export PROJECT_DESC="$project_desc"
export FEATURE_NAME="$feature_name"
export METRIC_NAMESPACE="$metric_namespace"
export DEFAULT_PORT="$default_port"

sed_replacement() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//&/\\&}"
  value="${value//|/\\|}"
  printf '%s' "$value"
}

project_name_sed="$(sed_replacement "$PROJECT_NAME")"
go_module_sed="$(sed_replacement "$GO_MODULE")"
project_desc_sed="$(sed_replacement "$PROJECT_DESC")"
feature_name_sed="$(sed_replacement "$FEATURE_NAME")"
metric_namespace_sed="$(sed_replacement "$METRIC_NAMESPACE")"
default_port_sed="$(sed_replacement "$DEFAULT_PORT")"

find "$target_dir" -type f -print0 | while IFS= read -r -d '' file; do
  sed -i.bak \
    -e "s|__PROJECT_NAME__|$project_name_sed|g" \
    -e "s|__GO_MODULE__|$go_module_sed|g" \
    -e "s|__PROJECT_DESC__|$project_desc_sed|g" \
    -e "s|__FEATURE_NAME__|$feature_name_sed|g" \
    -e "s|__METRIC_NAMESPACE__|$metric_namespace_sed|g" \
    -e "s|__DEFAULT_PORT__|$default_port_sed|g" \
    "$file"
  rm -f "$file.bak"
done

find "$target_dir" -depth -name '*__PROJECT_NAME__*' -print | while IFS= read -r path; do
  new_path="${path//__PROJECT_NAME__/$project_name}"
  mv "$path" "$new_path"
done

cat <<EOF
Rendered $project_name into $target_dir

Next:
  cd "$target_dir"
  go mod tidy
  make go-check
EOF
