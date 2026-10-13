#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

if [[ ${1:-} == --watch ]]; then
    rg --files src -g '*.go' -g go.mod -g Makefile | entr -r ./deploy.sh
    exit
fi

stack_name=${CDK_STACK_NAME:-$(jq -er 'if length == 1 then keys[0] else error("Set CDK_STACK_NAME to select a stack") end' outputs.json)}
function_name=$(jq -er --arg stack "$stack_name" '.[$stack].LambdaFunctionName' outputs.json)

make -C src
aws lambda update-function-code \
    --function-name "$function_name" \
    --zip-file fileb://src/function.zip \
    --query '[FunctionName,LastModified]' --output text --no-cli-pager
aws lambda wait function-updated-v2 --function-name "$function_name"
./scripts/smoke-test.sh
