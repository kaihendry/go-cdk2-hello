# Go HTTP server with TypeScript CDK

A standard Go `net/http` server on ARM64 Lambda, exposed through API Gateway HTTP
API. TypeScript defines the infrastructure; the Go application has no external
module dependencies. CDK builds the Linux binary locally and publishes a content
addressed asset, so no Docker build or manually managed S3 bucket is needed.

## Why Lambda Web Adapter?

[Lambda Web Functions](https://docs.aws.amazon.com/lambda/latest/dg/web-functions.html)
currently support only Node.js. To keep the application in Go, this project uses
[AWS Lambda Web Adapter](https://github.com/aws/aws-lambda-web-adapter) as an AWS
Lambda layer (ARM64 layer version 30), with the `provided.al2023` runtime. This
removes `apex/gateway` and all Lambda-specific code from the application. It is a
regular Lambda Event Function; it does not provide Web Functions' active CPU
billing or automatic multi-region deployment.

The adapter forwards requests to the Go server on port 8080. The executable is
named `bootstrap`; this custom-runtime setup does not need an exec wrapper.
API Gateway forwards all methods and paths to the server. Responses include the
build's Git revision in `X-Version` and in `CICD World ... <version>`.

## Local development and checks

Install Node.js 24+, npm, a supported Go release (at least 1.24), make, zip, the AWS
CLI, and jq. The optional watch shortcut also needs ripgrep and entr.

```sh
npm ci
make test
cd src
PORT=3000 go run .
```

The default port is 8080. A local `go run` reports version `dev`; deployment builds
embed `git describe --always --dirty --abbrev=7`. Set `VERSION` to override it.

## Deploy

```sh
export AWS_PROFILE=mine
aws sso login
npm ci
# Only needed once per account/region if CDK isn't already bootstrapped:
npm run cdk -- bootstrap
make deploy
```

CDK uses the selected profile's account and region. `AWS_REGION` can override the
region. Without domain configuration, the stack is `go-cdk2-hello` and the public
HTTPS URL is written to `outputs.json`. `make deploy` verifies the body and version
for GET, POST, PUT, DELETE, PATCH, OPTIONS, and HEAD, on both root and nested paths.
The endpoint is public, just like the original hello-world API.

Use `CDK_STACK_NAME` for an independent test stack. Custom domains are optional;
set both `DOMAIN` and `AWSCERT` (an ACM certificate in the deployment region):

```sh
export DOMAIN=hello.example.com
export AWSCERT=arn:aws:acm:eu-west-2:123456789012:certificate/your-certificate
make deploy
```

The default stack name then becomes `hello-example-com`. Point its DNS CNAME to
`APIGatewayDomain` from `outputs.json`. The GitHub workflow retains the existing
staging/production domains, deployment role, certificate, and Cloudflare secrets.
It runs checks before deploying staging, then deploys production after staging
passes. DNS is updated only when the CNAME target changes. Dependabot checks npm,
Go modules, and GitHub Actions weekly. Update the workflow settings when using
another AWS account or domain.

## Fast code updates

After the initial CDK deployment, keep the same AWS profile and region:

```sh
./deploy.sh
# Or rebuild and update whenever an existing source file changes:
./deploy.sh --watch
```

This uploads the Go ZIP directly with `aws lambda update-function-code`, waits for
the update, and runs the same HTTP smoke checks. It bypasses CloudFormation for
code-only changes. Use `make deploy` for infrastructure or layer changes; CDK will
publish the current source again as an asset. If `outputs.json` has multiple
stacks, select one with `CDK_STACK_NAME`.

To remove a test deployment, use the same stack name, profile, and region:

```sh
npm run cdk -- destroy
```
