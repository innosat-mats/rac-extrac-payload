import os
from aws_cdk import (
    Duration,
    Stack,
    aws_events as events,
    aws_events_targets as targets,
    aws_lambda as lambda_,
    aws_s3 as s3,
)
from constructs import Construct


RAC_INPUT_BUCKET = "mats-l0-rac"
RAC_OUTPUT_BUCKET = "mats-l0-artifacts"
RETENTION_PERIOD = 14 * 24 * 3600  # 14 days [s]


class RacLambdaStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        input_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacInputBucket",
            RAC_INPUT_BUCKET,
        )
        output_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacOutputBucket",
            RAC_OUTPUT_BUCKET,
        )

        rac_lambda = lambda_.Function(
            self,
            "rac-lambda",
            code=lambda_.InlineCode.from_asset("./handler"),
            handler="index.handler",
            timeout=Duration.seconds(900),
            runtime=lambda_.Runtime.PYTHON_3_9,
            environment={
                "RAC_OUTPUT_BUCKET": RAC_OUTPUT_BUCKET,
            },
        )

        input_bucket.add_event_notification()

        rule = events.Rule(
            self,
            "RacLambdaRule",
            schedule=events.Schedule.rate(Duration.hours(6)),
        )
        
        rule.add_target(targets.LambdaFunction(rac_lambda))

        input_bucket.grant_read(rac_lambda)
        output_bucket.grant_put(rac_lambda)
