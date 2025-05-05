package main

import (
	"os"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscertificatemanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

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
		"DOMAIN": jsii.String(os.Getenv("DOMAIN")),
	}

	codeBucketName := jsii.String("hendry-lambdas")
	codeBucket := awss3.Bucket_FromBucketName(stack, jsii.String("ImportedLambdaCodeBucket"), codeBucketName)
	codeObjectKey := jsii.String("go-cdk2-hello/function.zip")

	goURLFunction := awslambda.NewFunction(stack, jsii.String("go-function-al"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromBucket(codeBucket, codeObjectKey, nil),
		Environment:  funcEnvVar,
		Architecture: awslambda.Architecture_ARM_64(),
	})

	goURLFunction.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("s3:GetObject")},
		Resources: &[]*string{codeBucket.ArnForObjects(codeObjectKey)},
	}))

	awscdk.NewCfnOutput(stack, jsii.String("LambdaFunctionName"), &awscdk.CfnOutputProps{
		Value:       goURLFunction.FunctionName(),
		Description: jsii.String("Lambda function name"),
	})

	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String(domainName), &awsapigatewayv2.HttpApiProps{
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
		jsii.String("APIEndpoint"),
		&awscdk.CfnOutputProps{Value: httpApi.Url(), Description: jsii.String("API Gateway endpoint")},
	)

	awscdk.NewCfnOutput(
		stack,
		jsii.String("APIGatewayDomain"),
		&awscdk.CfnOutputProps{Value: dn.RegionalDomainName(), Description: jsii.String("API Gateway domain name")},
	)

	awscdk.NewCfnOutput(
		stack,
		jsii.String("CustomDomain"),
		&awscdk.CfnOutputProps{Value: jsii.String(domainName), Description: jsii.String("Custom domain name")},
	)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)
	stackName := strings.ReplaceAll(domainName, ".", "-")
	if stackName == "" {
		stackName = "go-cdk2-hello"
	}
	NewGStack(app, stackName, &GStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region:  jsii.String(os.Getenv("AWS_REGION")),
		Account: jsii.String("407461997746"),
	}
}
