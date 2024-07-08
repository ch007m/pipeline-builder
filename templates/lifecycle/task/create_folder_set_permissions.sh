#!/usr/bin/env bash
set -e

# TODO: To be reviewed and aligned with Shipwright ClusterBuildStrategy
if [[ "$(workspaces.cache.bound)" == "true" ]]; then
  echo "> Setting permissions on '$(workspaces.cache.path)'..."
  chown -R "$(params.CNB_USER_ID):$(params.CNB_GROUP_ID)" "$(workspaces.cache.path)"
fi

echo "Creating .docker folder"
mkdir -p "/tekton/home/.docker"

for path in "/tekton/home" "/tekton/home/.docker" "/layers" "$(workspaces.source.path)"; do
  echo "> Setting permissions on '$path'..."
  chown -R "$(params.CNB_USER_ID):$(params.CNB_GROUP_ID)" "$path"
done

echo "> Parsing additional configuration..."
parsing_flag=""
envs=()
for arg in "$@"; do
    if [[ "$arg" == "--env-vars" ]]; then
        echo "-> Parsing env variables..."
        parsing_flag="env-vars"
    elif [[ "$parsing_flag" == "env-vars" ]]; then
        envs+=("$arg")
    fi
done

echo "> Processing any environment variables..."
ENV_DIR="/platform/env"

echo "--> Creating 'env' directory: $ENV_DIR"
mkdir -p "$ENV_DIR"

for env in "${envs[@]}"; do
    IFS='=' read -r key value string <<< "$env"
    if [[ "$key" != "" && "$value" != "" ]]; then
        path="${ENV_DIR}/${key}"
        echo "--> Writing ${path}..."
        echo -n "$value" > "$path"
    fi
done