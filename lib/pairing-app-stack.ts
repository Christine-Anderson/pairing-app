import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as path from "path";

export class PairingAppStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const table = new dynamodb.Table(this, "MyGroupTable", {
            partitionKey: { 
                name: "groupId",
                type: dynamodb.AttributeType.STRING
            },
        });

        const myFunction = new lambda.Function(this, "myLambdaFunction", {
            runtime: lambda.Runtime.PROVIDED_AL2023,
            handler: "main",
            code: lambda.Code.fromAsset(path.join(__dirname, "../lambda/function.zip")),
        });

        table.grantReadWriteData(myFunction)

        const api = new apigateway.RestApi(this, "myApiGateway", {
            defaultCorsPreflightOptions: {
                allowHeaders: ["Content-type", "Authorization"],
                allowMethods: ["GET", "POST"],
                allowOrigins: ["*"]
            },
        })

        const integration = new apigateway.LambdaIntegration(myFunction)

        const createGroupResource = api.root.addResource("create-group");
        createGroupResource.addMethod("POST", integration);

        const joinGroupResource = api.root.addResource("join-group");
        joinGroupResource.addMethod("POST", integration);

        const groupDetailsResource = api.root.addResource("group-details");
        groupDetailsResource.addResource("{groupId}");
        groupDetailsResource.addMethod("POST", integration);

        const performMatchingResource = api.root.addResource("match");
        performMatchingResource.addMethod("POST", integration);
    }
}
