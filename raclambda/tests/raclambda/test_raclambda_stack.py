import pytest  # type: ignore

from aws_cdk import App
from aws_cdk.assertions import Match, Template

from raclambda.raclambda_stack import RacLambdaStack


@pytest.fixture
def template():
    app = App()

    stack = RacLambdaStack(
        app,
        "raclambda",
        "input-bucket",
        "output-bucket",
        "test-project",
    )

    return Template.from_stack(stack)


class TestRacLambdaStack:

    def test_has_lambda_policy(self, template: Template):
        template.has_resource_properties(
            "AWS::IAM::Policy",
            {
                "PolicyDocument": {
                    "Statement": [
                        {
                            "Action": [
                                "s3:GetObject*",
                                "s3:GetBucket*",
                                "s3:List*"
                            ],
                            "Effect": "Allow",
                            "Resource": [
                                {
                                    "Fn::Join": [
                                        "",
                                        [
                                            "arn:",
                                            {
                                                "Ref": "AWS::Partition"
                                            },
                                            ":s3:::input-bucket"
                                        ]
                                    ]
                                },
                                {
                                    "Fn::Join": [
                                        "",
                                        [
                                            "arn:",
                                            {
                                                "Ref": "AWS::Partition"
                                            },
                                            ":s3:::input-bucket/*"
                                        ]
                                    ]
                                }
                            ]
                        },
                        {
                            "Action": [
                                "s3:PutObject",
                                "s3:PutObjectLegalHold",
                                "s3:PutObjectRetention",
                                "s3:PutObjectTagging",
                                "s3:PutObjectVersionTagging",
                                "s3:Abort*"
                            ],
                            "Effect": "Allow",
                            "Resource": {
                                "Fn::Join": [
                                    "",
                                    [
                                        "arn:",
                                        {
                                            "Ref": "AWS::Partition"
                                        },
                                        ":s3:::output-bucket/*"
                                    ]
                                ]
                            }
                        },
                        {
                            "Action": [
                                "sqs:ReceiveMessage",
                                "sqs:ChangeMessageVisibility",
                                "sqs:GetQueueUrl",
                                "sqs:DeleteMessage",
                                "sqs:GetQueueAttributes"
                            ],
                            "Effect": "Allow",
                            "Resource": {
                                "Fn::GetAtt": [
                                    "RacQueue12CAA348",
                                    "Arn"
                                ]
                            }
                        }
                    ],
                    "Version": "2012-10-17"
                },
                "PolicyName": "raclambdaServiceRoleDefaultPolicy0908116C",
                "Roles": [
                    {
                        "Ref": "raclambdaServiceRole61297EF8"
                    }
                ]
            }
        )

    def test_has_lambda_function(self, template: Template):
        template.has_resource_properties(
            "AWS::Lambda::Function",
            {
                "Code": {
                    "S3Bucket": {
                        "Fn::Sub": Match.string_like_regexp(".*")
                    },
                    "S3Key": Match.string_like_regexp(".*\\.zip")
                },
                "Role": {
                    "Fn::GetAtt": [
                        "raclambdaServiceRole61297EF8",
                        "Arn"
                    ]
                },
                "Environment": {
                    "Variables": {
                        "RAC_INPUT_BUCKET": "input-bucket",
                        "RAC_QUEUE": {
                            "Fn::GetAtt": [
                                "RacQueue12CAA348",
                                "QueueName"
                            ]
                        },
                        "RAC_PROJECT": "test-project"
                    }
                },
                "Handler": "raclambda_handler.handler",
                "Runtime": "python3.9",
                "Timeout": 300
            },
        )

    def test_has_lambda_event(self, template: Template):
        template.has_resource_properties(
            "AWS::Events::Rule",
            {
                "ScheduleExpression": "rate(6 hours)",
                "State": "ENABLED",
                "Targets": [
                    {
                        "Arn": {
                            "Fn::GetAtt": [
                                "raclambda99ECE2E1",
                                "Arn"
                            ]
                        },
                        "Id": "Target0"
                    }
                ]
            }
        )

        template.has_resource_properties(
            "AWS::Lambda::Permission",
            {
                "Action": "lambda:InvokeFunction",
                "Principal": "events.amazonaws.com"
            }
        )
