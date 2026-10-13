import { App } from "aws-cdk-lib";
import { HelloStack } from "./stack";

const app = new App();
const domainName = process.env.DOMAIN || undefined;
const certificateArn = process.env.AWSCERT || undefined;

new HelloStack(
  app,
  process.env.CDK_STACK_NAME || domainName?.replaceAll(".", "-") || "go-cdk2-hello",
  {
    env: {
      account: process.env.CDK_DEFAULT_ACCOUNT,
      region: process.env.AWS_REGION || process.env.CDK_DEFAULT_REGION,
    },
    domainName,
    certificateArn,
  },
);

app.synth();
