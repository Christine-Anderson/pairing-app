import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as path from 'path';

export class PairingAppStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Define the Lambda function resource
        const myFunction = new lambda.Function(this, "HelloWorldFunction", {
                runtime: lambda.Runtime.PROVIDED_AL2023,
                handler: "main",
                code: lambda.Code.fromAsset(path.join(__dirname, '../lambda/function.zip')),
            });

        // Define the Lambda function URL resource
        const myFunctionUrl = myFunction.addFunctionUrl({
            authType: lambda.FunctionUrlAuthType.NONE,
        });

        // Define a CloudFormation output for your URL
        new cdk.CfnOutput(this, "myFunctionUrlOutput", {
            value: myFunctionUrl.url,
        })
    }
}
