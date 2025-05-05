#!/bin/bash

# Exit on error
set -e

# Check if entr is installed
if ! command -v entr >/dev/null 2>&1; then
    echo "entr is required but not installed. Please install it first."
    echo "On macOS: brew install entr"
    echo "On Linux: sudo apt-get install entr"
    exit 1
fi

# Check if aws cli is installed
if ! command -v aws >/dev/null 2>&1; then
    echo "AWS CLI is required but not installed. Please install it first."
    echo "See: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html"
    exit 1
fi

# Check if jq is installed
if ! command -v jq >/dev/null 2>&1; then
    echo "jq is required but not installed. Please install it first."
    echo "On macOS: brew install jq"
    echo "On Linux: sudo apt-get install jq"
    exit 1
fi

# Check if CDK_STACK_NAME is set
if [ -z "$CDK_STACK_NAME" ]; then
    echo "Error: CDK_STACK_NAME environment variable is not set."
    exit 1
fi
echo "Using CDK Stack Name: $CDK_STACK_NAME"

# Get values from outputs.json
if [ ! -f outputs.json ]; then
    echo "outputs.json not found. Please run 'cdk deploy --outputs-file outputs.json' first."
    exit 1
fi

# Get the Lambda function name from outputs.json using the stack name
LAMBDA_FUNCTION_NAME=$(jq -r '."'$CDK_STACK_NAME'".LambdaFunctionName' outputs.json)
if [ -z "$LAMBDA_FUNCTION_NAME" ] || [ "$LAMBDA_FUNCTION_NAME" = "null" ]; then
    echo "Could not find LambdaFunctionName for stack '$CDK_STACK_NAME' in outputs.json."
    exit 1
fi

# Define the target S3 bucket and key (matching g.go)
BUCKET_NAME="hendry-lambdas"
S3_KEY="go-cdk2-hello/function.zip"
S3_DESTINATION="s3://${BUCKET_NAME}/${S3_KEY}"

echo "Found Lambda function: $LAMBDA_FUNCTION_NAME"
echo "Target S3 location: $S3_DESTINATION"
echo "Watching for changes in src/main.go..."

# The deploy function that will be run when changes are detected
deploy() {
    echo "Change detected in main.go, rebuilding..."

    # Build the project (creates src/function.zip)
    if ! make -C src/; then
        echo "Build failed!"
        return 1
    fi
    echo "Build successful."

    # Upload to S3
    echo "Uploading function.zip to $S3_DESTINATION..."
    if ! aws s3 cp src/function.zip "$S3_DESTINATION"; then
        echo "S3 upload failed!"
        return 1
    fi
    echo "S3 upload successful."

    # Update Lambda function code using AWS CLI
    echo "Updating Lambda function code for $LAMBDA_FUNCTION_NAME..."
    if ! aws lambda update-function-code --function-name "$LAMBDA_FUNCTION_NAME" --s3-bucket "$BUCKET_NAME" --s3-key "$S3_KEY" --no-cli-pager; then
        echo "Lambda update failed!"
        return 1
    fi

    # Wait for the update to complete
    echo "Waiting for Lambda function update to complete..."
    if ! aws lambda wait function-updated --function-name "$LAMBDA_FUNCTION_NAME"; then
        echo "Lambda function update wait failed!"
        return 1
    fi

    echo "Lambda function update complete!"
}

# Watch src/main.go and run the deploy function on changes
# Pass required variables to the subshell environment for entr
export LAMBDA_FUNCTION_NAME
export BUCKET_NAME
export S3_KEY
export S3_DESTINATION

# If run directly with `--watch`, start the watch loop
if [[ "$1" == "--watch" ]]; then
    # Use a here-string to pass file list to entr
    find src -type f \
        ! -name 'bootstrap' \
        ! -name '*.zip' |
        entr -r "$0"
    exit
fi

# Otherwise, just run deploy
deploy
# curl APIEndpoint
echo "Invoking the updated function..."
API_ENDPOINT=$(jq -r '."'$CDK_STACK_NAME'".APIEndpoint' outputs.json)
if [ -z "$API_ENDPOINT" ] || [ "$API_ENDPOINT" = "null" ]; then
    echo "Could not find APIEndpoint for stack '$CDK_STACK_NAME' in outputs.json."
    exit 1
fi
curl "$API_ENDPOINT"
