# Go Hello World using AWS CDK 2

## Environment variables in .github/workflows

For example using DNS validation with wildcard *.dabase.com

    AWSCERT="arn:aws:acm:ap-southeast-1:407461997746:certificate/87b0fd84-fb44-4782-b7eb-d9c7f8714908"
    DOMAIN="hello.dabase.com"

CI/CD deployment; [you need to adjust _role-to-assume_ for the workflow to work](https://youtu.be/WKzVqFsOBSE), once setup you avoid the need to setup AWS_SECRET_ACCESS_KEY credentials.

# To deploy to the cloud

    npx aws-cdk@2.x deploy

# To develop locally

    cd src
    go get github.com/codegangsta/gin
    gin

# Why?

There are many ways to deploy a Go application to the AWS Cloud. I've explored them all.

## EC2

Copy across the Go binary and put it behind an ALB or Caddy.

Awkward and isn't serverless.

## Make a Docker image and deploy with Kubernetes

Have you lost your mind? This is incredibly complex and expensive way to deploy a Go application to the Cloud.

Good luck to you and your team.

And it's not serverless.

## AWS Serverless Application Model (SAM)

https://github.com/kaihendry/aws-sam-gateway-example

An efficient usage of AWS native Cloudformation that requires Python tooling such as https://github.com/aws/aws-sam-cli

## Apex Up

Easiest though proprietary https://apex.sh/up/

## Terraform

Very awkward and slow to deploy and requires extra tooling.

## AWS CDK

Keep everything in Go, including the Infrastructure as Code.

Warning experimental APIs: https://pkg.go.dev/github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2
