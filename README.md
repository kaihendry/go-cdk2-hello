# Go HTTP server with TypeScript CDK

A dependency-free Go HTTP server on ARM64 Lambda, behind API Gateway.
[Lambda Web Functions](https://docs.aws.amazon.com/lambda/latest/dg/web-functions.html)
currently support only Node.js, so this uses the
[Lambda Web Adapter](https://github.com/aws/aws-lambda-web-adapter) layer instead.
It listens on port 8080 and returns its Git revision in the body and `X-Version`.

[SLOC history](https://kaihendry.github.io/go-cdk2-hello/) is published by
[LessHero](https://github.com/kaihendry/lesshero) on every push to `main`.
It counts tracked lines across Git history, excluding dependency lockfiles.

## Develop

Use Node.js 24 (`.node-version`), the pnpm version pinned in `package.json`,
a supported Go release, make, zip, the AWS CLI, and jq.

```sh
pnpm install --frozen-lockfile
make test
cd src && PORT=3000 go run .
```

## Deploy

```sh
export AWS_PROFILE=mine
aws sso login
# Once per account/region, if needed:
pnpm cdk bootstrap
make deploy
curl -fsS "$(jq -r '.[].APIEndpoint' outputs.json)"
```

`make deploy` builds `src/function.zip`; CDK uploads it as an asset. No Docker or
manually managed bucket is needed. CDK uses the profile's account and region;
`AWS_REGION` overrides the region. `make synth` builds and previews the template.

Set `CDK_STACK_NAME` for an independent test stack (default: `go-cdk2-hello`).
For a custom domain, set both `DOMAIN` and `AWSCERT` (an ACM certificate in the
same region). The default stack name then replaces dots in `DOMAIN` with hyphens.
Point the domain's CNAME to `APIGatewayDomain` from `outputs.json`.

For fast code-only updates, run `./deploy.sh`, or `./deploy.sh --watch` with
ripgrep and entr installed. This uploads directly to Lambda, waits, and requests
the endpoint. Keep the same AWS profile/region and select `CDK_STACK_NAME` if
`outputs.json` contains multiple stacks. Set `VERSION` to override the embedded
Git revision. Remove a deployment with `pnpm cdk destroy` using the same settings.

## CI

Checks and staging run in parallel. Staging builds and uploads the Lambda ZIP;
production downloads that exact artifact after both checks and staging succeed.
The smoke check requests the custom domain and asserts the deployed Git revision.
Daily runs only check the code and vulnerabilities; deployments run on pushes to
`main` and manual dispatches.

The deployment workflow contains this project's AWS role, certificate, and
staging/production domains. DNS changes use `CLOUDFLARE_API_TOKEN` and
`CLOUDFLARE_ZONE`; no API call is made if the CNAME already matches.
Dependabot checks pnpm, Go modules, and Actions weekly, keeping Node types on v24.
