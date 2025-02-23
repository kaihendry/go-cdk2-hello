package main

import (
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/aws-cdk-go/awscdk/v2"

	"github.com/aws/constructs-go/constructs/v10"

	"github.com/aws/jsii-runtime-go"
)

var (
	certArn    = os.Getenv("AWSCERT")
	domainName = os.Getenv("DOMAIN")
)

type GStackProps struct {
	awscdk.StackProps
}

func NewGStack(scope constructs.Construct, id string, props *GStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	dn := awsapigatewayv2.NewDomainName(stack, jsii.String("DN"), &awsapigatewayv2.DomainNameProps{
		DomainName:  jsii.String(domainName),
		Certificate: awscertificatemanager.Certificate_FromCertificateArn(stack, jsii.String("Cert"), jsii.String(certArn)),
	})

	funcEnvVar := &map[string]*string{
		"VERSION": jsii.String(os.Getenv("VERSION")),
	}

	goURLFunction := awslambda.NewFunction(stack, jsii.String("go-function"), &awslambda.FunctionProps{
		Runtime:     awslambda.Runtime_PROVIDED_AL2023(),
		Handler:     jsii.String("bootstrap"),
		Code:        awslambda.AssetCode_FromAsset(jsii.String("src/function.zip"), nil),
		Environment: funcEnvVar,
	})

	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("MyHttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("MyHttpApi"),
		DefaultDomainMapping: &awsapigatewayv2.DomainMappingOptions{
			DomainName: dn,
		},
	})

	lambdaIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(jsii.String("LambdaIntegration"), goURLFunction, nil)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.String("/"),
		Integration: lambdaIntegration,
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
			awsapigatewayv2.HttpMethod_POST,
			awsapigatewayv2.HttpMethod_PUT,
			awsapigatewayv2.HttpMethod_DELETE,
			awsapigatewayv2.HttpMethod_PATCH,
			awsapigatewayv2.HttpMethod_OPTIONS,
			awsapigatewayv2.HttpMethod_HEAD,
		},
	})

	awscdk.NewCfnOutput(
		stack,
		jsii.String("API Endpoint"),
		&awscdk.CfnOutputProps{Value: httpApi.Url(), Description: jsii.String("API Gateway endpoint")},
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	NewGStack(app, strings.ReplaceAll(domainName, ".", "-"), &GStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
