import os
from aws_cdk import (
    Duration,
    Fn,
    Stack,
    aws_lambda_event_sources as sources,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_sqs as sqs,
)
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
        slask_bucket_name: str,
        project_name: str,
        queue_arn_export_name: str,
        lambda_timeout: Duration = Duration.seconds(300),
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
        # TODO: create with retention rule
        slask_bucket = s3.Bucket.from_bucket_name(
            self,
            "RacSlaskBucket",
            slask_bucket_name,
        )
        rac_queue = sqs.Queue.from_queue_arn(
            self,
            "RacQueue",
            Fn.import_value(queue_arn_export_name)
        )

        rac_lambda = lambda_.Function(
            self,
            "rac-lambda",
            code=lambda_.InlineCode.from_asset("./raclambda/handler"),
            handler="raclambda_handler.handler",
            timeout=lambda_timeout,
            runtime=lambda_.Runtime.PYTHON_3_9,
            environment={
                "RAC_PROJECT": project_name,
                "RAC_SLASK": slask_bucket_name,
            },
        )

        rac_lambda.add_event_source(sources.SqsEventSource(
            rac_queue,
            batch_size=1,
        ))

        input_bucket.grant_read(rac_lambda)
        output_bucket.grant_put(rac_lambda)
        slask_bucket.grant_read_write(rac_lambda)
