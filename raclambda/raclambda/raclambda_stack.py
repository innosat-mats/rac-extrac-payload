from aws_cdk import (
    Duration,
    Stack,
    # aws_sqs as sqs,
    aws_events as events,
    aws_events_targets as targets,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_sqs as sqs,
)
from constructs import Construct


INPUT_BUCKET = "mats-l0-rac"
OUTPUT_BUCKET = "mats-l0-artifacts"
RAC_QUEUE_ARN = "arn:aws:sqs:for:eu-north-1:12345:racqueue1"


class RacLambdaStack(Stack):

    def __init__(self, scope: Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        input_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacInputBucket",
            INPUT_BUCKET,
        )
        output_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacOutputBucket",
            OUTPUT_BUCKET,
        )
        rac_queue = sqs.Queue.from_queue_arn(
            self,
            "RacQueue",
            RAC_QUEUE_ARN,
        )

        rac_lambda = lambda_.Function(
            self,
            "rac-lambda",
            code=lambda_.InlineCode.from_asset("./handler"),
            handler="index.handler",
            timeout=Duration.seconds(900),
            runtime=lambda_.Runtime.PYTHON_3_9,
            environment={
                "INPUT_BUCKET": INPUT_BUCKET,
                "OUTPUT_BUCKET": OUTPUT_BUCKET,
                "RAC_QUEUE_ARN": RAC_QUEUE_ARN,
            },
        )

        rule = events.Rule(
            self,
            "RacLambdaRule",
            schedule=events.Schedule.rate(Duration.hours(2)),
        )
        
        rule.add_target(targets.LambdaFunction(rac_lambda))

        input_bucket.grant_read(rac_lambda)
        output_bucket.grant_put(rac_lambda)
        rac_queue.grant_consume_messages(rac_lambda)
