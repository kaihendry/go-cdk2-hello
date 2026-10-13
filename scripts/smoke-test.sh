#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

stack_name=${CDK_STACK_NAME:-$(jq -er 'if length == 1 then keys[0] else error("Set CDK_STACK_NAME to select a stack") end' outputs.json)}
endpoint=${API_ENDPOINT:-$(jq -er --arg stack "$stack_name" '.[$stack].APIEndpoint' outputs.json)}
expected_version=${VERSION:-$(git describe --always --dirty --abbrev=7)}
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

for path in / '/nested/path?test=1'; do
    for method in GET POST PUT DELETE PATCH OPTIONS HEAD; do
        args=(--request "$method")
        if [[ $method == HEAD ]]; then args=(--head); fi
        curl --fail --silent --show-error --connect-timeout 10 --max-time 30 \
            --retry 5 --retry-delay 2 --retry-all-errors \
            "${args[@]}" --dump-header "$tmp/headers" --output "$tmp/body" \
            "${endpoint%/}$path"
        actual_version=$(awk 'tolower($1) == "x-version:" {gsub("\r", "", $2); print $2}' "$tmp/headers")
        if [[ $actual_version != "$expected_version" ]]; then
            echo "Unexpected version: '$actual_version' (expected '$expected_version')" >&2
            exit 1
        fi
        if [[ $method != HEAD && $(cat "$tmp/body") != "CICD World ... $expected_version" ]]; then
            echo "Unexpected body for $method $path" >&2
            cat "$tmp/body" >&2
            exit 1
        fi
        echo "OK $method $path ($actual_version)"
    done
done
