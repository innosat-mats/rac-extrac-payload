import os
from aws_cdk import (
    Duration,
    Stack,
    aws_events_targets as targets,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_s3_notifications as s3n,
    aws_sqs as sqs,
)
from aws_cdk.aws_events import Rule, Schedule
from constructs import Construct


RAC_PROJECT = os.environ.get("RAC_PROJECT", "mats-test-project")
RETENTION_PERIOD = 14 * 24 * 3600  # 14 days [s]


class RacLambdaStack(Stack):

    def __init__(
        self,
        scope: Construct,
        construct_id: str,
        input_bucket_name: str,
        output_bucket_name: str,
        project_name: str,
        lambda_timeout: Duration = Duration.seconds(300),
        lambda_schedule: Schedule = Schedule.rate(Duration.hours(6)),
        queue_retention: Duration = Duration.days(14),
        **kwargs,
    ) -> None:
        super().__init__(scope, construct_id, **kwargs)

        input_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacInputBucket",
            input_bucket_name,
        )
        output_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacOutputBucket",
            output_bucket_name,
        )
        rac_queue = sqs.Queue(
            self,
            "RacQueue",
            retention_period=queue_retention,
        )

        rac_lambda = lambda_.Function(
            self,
            "rac-lambda",
            code=lambda_.InlineCode.from_asset("./raclambda/handler"),
            handler="raclambda_handler.handler",
            timeout=lambda_timeout,
            runtime=lambda_.Runtime.PYTHON_3_9,
            environment={
                "RAC_INPUT_BUCKET": input_bucket_name,
                "RAC_QUEUE": rac_queue.queue_name,
                "RAC_PROJECT": project_name,
            },
        )

        rule = Rule(
            self,
            "RacLambdaRule",
            schedule=lambda_schedule,
        )
        rule.add_target(targets.LambdaFunction(rac_lambda))

        input_bucket.grant_read(rac_lambda)
        output_bucket.grant_put(rac_lambda)
        rac_queue.grant_consume_messages(rac_lambda)

        input_bucket.add_event_notification(
            s3.EventType.OBJECT_CREATED,
            s3n.SqsDestination(rac_queue),
        )
