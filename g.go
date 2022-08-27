package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/aws-cdk-go/awscdk/v2"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"

	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"

	"github.com/aws/jsii-runtime-go"
)

type GStackProps struct {
	awscdk.StackProps
}

const accessURLFunctionDirectory = "access"

func NewGStack(scope constructs.Construct, id string, props *GStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	certArn := "arn:aws:acm:ap-southeast-1:407461997746:certificate/87b0fd84-fb44-4782-b7eb-d9c7f8714908"
	domainName := "hello.dabase.com"

	dn := awscdkapigatewayv2alpha.NewDomainName(stack, jsii.String("DN"), &awscdkapigatewayv2alpha.DomainNameProps{
		DomainName:  jsii.String(domainName),
		Certificate: awscertificatemanager.Certificate_FromCertificateArn(stack, jsii.String("Cert"), jsii.String(certArn)),
	})

	acccessURLFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("access-url-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime: awslambda.Runtime_GO_1_X(),
			Entry:   jsii.String(accessURLFunctionDirectory)})

	accessFunctionIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("access-function-integration"), acccessURLFunction, nil)

	httpApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("HttpProxyProdApi"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("HttpProxyProdApi"),
		DefaultDomainMapping: &awscdkapigatewayv2alpha.DomainMappingOptions{
			DomainName: dn,
		},
		DefaultIntegration: accessFunctionIntg,
	})

	httpApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: accessFunctionIntg})

	awscdk.NewCfnOutput(stack, jsii.String("output"), &awscdk.CfnOutputProps{Value: httpApi.Url(), Description: jsii.String("API Gateway endpoint")})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewGStack(app, "hello-go", &GStackProps{
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
