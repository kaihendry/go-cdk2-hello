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

# Get values from outputs.json
if [ ! -f outputs.json ]; then
    echo "outputs.json not found. Please run 'cdk deploy' first."
    exit 1
fi

# Get the Lambda function name from outputs.json
LAMBDA_FUNCTION_NAME=$(jq -r '."stghello-dabase-com".LambdaFunctionName' outputs.json)
if [ -z "$LAMBDA_FUNCTION_NAME" ] || [ "$LAMBDA_FUNCTION_NAME" = "null" ]; then
    echo "Could not find LambdaFunctionName in outputs.json. Please run 'cdk deploy' first."
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
    if ! aws lambda update-function-code --function-name "$LAMBDA_FUNCTION_NAME" --s3-bucket "$BUCKET_NAME" --s3-key "$S3_KEY"; then
        echo "Lambda update failed!"
        return 1
    fi

    echo "Lambda function updated successfully!"
}

# Watch src/main.go and run the deploy function on changes
# Pass required variables to the subshell environment for entr
export LAMBDA_FUNCTION_NAME
export BUCKET_NAME
export S3_KEY
export S3_DESTINATION
echo src/main.go | entr -n -r bash -c "$(declare -f deploy); deploy"
