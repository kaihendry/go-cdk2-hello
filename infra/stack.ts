import * as path from "node:path";
import {
  CfnOutput,
  Duration,
  Stack,
  StackProps,
  aws_apigatewayv2 as apigateway,
  aws_apigatewayv2_integrations as integrations,
  aws_certificatemanager as acm,
  aws_lambda as lambda,
} from "aws-cdk-lib";
import { Construct } from "constructs";

export interface HelloStackProps extends StackProps {
  domainName?: string;
  certificateArn?: string;
  code?: lambda.Code;
}

export class HelloStack extends Stack {
  constructor(scope: Construct, id: string, props: HelloStackProps = {}) {
    super(scope, id, props);

    if (Boolean(props.domainName) !== Boolean(props.certificateArn)) {
      throw new Error("Set both DOMAIN and AWSCERT, or leave both unset.");
    }

    const adapter = lambda.LayerVersion.fromLayerVersionArn(
      this,
      "WebAdapter",
      `arn:${this.partition}:lambda:${this.region}:753240598075:layer:LambdaAdapterLayerArm64:30`,
    );

    // Keep the construct ID so existing deployments retain their function.
    const fn = new lambda.Function(this, "go-function-al", {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      architecture: lambda.Architecture.ARM_64,
      handler: "bootstrap",
      code: props.code ?? lambda.Code.fromAsset(path.join(__dirname, "..", "src", "function.zip")),
      timeout: Duration.seconds(10),
      layers: [adapter],
      environment: {
        PORT: "8080",
        AWS_LWA_PORT: "8080",
        DOMAIN: props.domainName ?? "",
      },
    });

    const domain = props.domainName && props.certificateArn
      ? new apigateway.DomainName(this, "DN", {
          domainName: props.domainName,
          certificate: acm.Certificate.fromCertificateArn(this, "Cert", props.certificateArn),
        })
      : undefined;

    const api = new apigateway.HttpApi(this, props.domainName ?? "HttpApi", {
      defaultDomainMapping: domain ? { domainName: domain } : undefined,
      defaultIntegration: new integrations.HttpLambdaIntegration("LambdaIntegration", fn),
    });

    new CfnOutput(this, "LambdaFunctionName", { value: fn.functionName });
    new CfnOutput(this, "APIEndpoint", { value: api.apiEndpoint });
    if (domain) {
      new CfnOutput(this, "APIGatewayDomain", { value: domain.regionalDomainName });
      new CfnOutput(this, "CustomDomain", { value: props.domainName! });
    }
  }
}
