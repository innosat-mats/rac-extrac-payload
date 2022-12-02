from aws_cdk import (
    Duration,
    Fn,
    Size,
    Stack,
    aws_lambda_event_sources as sources,
    aws_lambda as lambda_,
    aws_s3 as s3,
    aws_sqs as sqs,
    aws_iam as iam,
)
from constructs import Construct


RETENTION_PERIOD = 14 * 24 * 3600  # 14 days [s]


class RacLambdaStack(Stack):

    def __init__(
        self,
        scope: Construct,
        construct_id: str,
        input_bucket_name: str,
        output_bucket_name: str,
        queue_arn_export_name: str,
        config_ssm_name: str,
        rclone_arn: str,
        lambda_timeout: Duration = Duration.seconds(300),
        dregs_expiration: Duration = Duration.days(7),
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

        dregs_bucket = s3.Bucket(
            self,
            "RacDregsBucket",
        )
        dregs_bucket.add_lifecycle_rule(expiration=dregs_expiration)

        rac_queue = sqs.Queue.from_queue_arn(
            self,
            "RacQueue",
            Fn.import_value(queue_arn_export_name)
        )

        rclone_layer = lambda_.LayerVersion.from_layer_version_arn(
            self,
            "RCloneLayer",
            rclone_arn,
        )

        rac_lambda = lambda_.Function(
            self,
            "rac-lambda",
            code=lambda_.InlineCode.from_asset("./raclambda/handler"),
            handler="raclambda_handler.handler",
            timeout=lambda_timeout,
            architecture=lambda_.Architecture.X86_64,
            runtime=lambda_.Runtime.PYTHON_3_9,
            memory_size=1024,
            ephemeral_storage_size=Size.mebibytes(1024),
            environment={
                "RAC_DREGS": dregs_bucket.bucket_name,
                "RAC_OUTPUT": output_bucket.bucket_name,
                "RCLONE_CONFIG_SSM_NAME": config_ssm_name,
            },
            layers=[rclone_layer],
        )

        rac_lambda.add_event_source(sources.SqsEventSource(
            rac_queue,
            batch_size=1,
        ))

        rac_lambda.add_to_role_policy(iam.PolicyStatement(
            effect=iam.Effect.ALLOW,
            actions=["ssm:GetParameter"],
            resources=[f"arn:aws:ssm:*:*:parameter{config_ssm_name}"]
        ))

        input_bucket.grant_read(rac_lambda)
        output_bucket.grant_read_write(rac_lambda)
        dregs_bucket.grant_read_write(rac_lambda)
