import assert from "node:assert/strict";
import { test } from "node:test";
import { App, aws_lambda as lambda } from "aws-cdk-lib";
import { Match, Template } from "aws-cdk-lib/assertions";
import { HelloStack } from "./stack";

const env = { account: "123456789012", region: "eu-west-2" };
// Inject deployment parameters so infrastructure tests do not compile Go.
const code = () => lambda.Code.fromCfnParameters();

test("runs the plain Go HTTP server through the ARM64 Web Adapter", () => {
  const template = Template.fromStack(new HelloStack(new App(), "Test", { env, code: code() }));
  template.hasResourceProperties("AWS::Lambda::Function", {
    Runtime: "provided.al2023",
    Handler: "bootstrap",
    Architectures: ["arm64"],
    Layers: [{ "Fn::Join": ["", [
      "arn:", { Ref: "AWS::Partition" },
      ":lambda:eu-west-2:753240598075:layer:LambdaAdapterLayerArm64:30",
    ]] }],
    Environment: { Variables: { PORT: "8080", AWS_LWA_PORT: "8080", DOMAIN: "" } },
  });
  template.hasResourceProperties("AWS::ApiGatewayV2::Route", { RouteKey: "$default" });
  template.hasResourceProperties("AWS::ApiGatewayV2::Integration", { PayloadFormatVersion: "2.0" });
  template.resourceCountIs("AWS::ApiGatewayV2::DomainName", 0);
  template.resourceCountIs("AWS::S3::Bucket", 0);
  template.hasOutput("APIEndpoint", { Value: Match.anyValue() });
  template.hasOutput("LambdaFunctionName", { Value: Match.anyValue() });
});

test("retains custom domain support and the existing Lambda logical ID", () => {
  const template = Template.fromStack(new HelloStack(new App(), "stghello-dabase-com", {
    env,
    code: code(),
    domainName: "stghello.dabase.com",
    certificateArn: "arn:aws:acm:eu-west-2:123456789012:certificate/example",
  }));
  template.hasResourceProperties("AWS::ApiGatewayV2::DomainName", {
    DomainName: "stghello.dabase.com",
  });
  template.resourceCountIs("AWS::ApiGatewayV2::ApiMapping", 1);
  assert.equal(Object.keys(template.findResources("AWS::ApiGatewayV2::Api"))[0], "stghellodabasecom73858662");
  assert.equal(Object.keys(template.findResources("AWS::ApiGatewayV2::DomainName"))[0], "DNFDC76583");
  template.hasOutput("CustomDomain", { Value: "stghello.dabase.com" });
  assert.equal(Object.keys(template.findResources("AWS::Lambda::Function"))[0], "gofunctionalC3901643");
});

test("rejects an incomplete custom domain configuration", () => {
  for (const config of [{ domainName: "hello.example.com" }, { certificateArn: "arn:example" }]) {
    assert.throws(
      () => new HelloStack(new App(), "Test", { env, code: code(), ...config }),
      /Set both DOMAIN and AWSCERT/,
    );
  }
});
