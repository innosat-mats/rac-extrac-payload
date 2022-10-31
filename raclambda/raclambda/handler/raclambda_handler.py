import json
import os
import subprocess
from typing import Any, List

import boto3


MIN_AGE = 6 * 3600  # 6 hours [s]

SQSQueue = Any
SQSMessage = Any


class NoNewFiles(Exception):
    pass


def get_env_or_raise(envvar: str) -> str:
    if (val := os.environ.get(envvar)) is None:
        raise ValueError(f"'{envvar}' not found in env")
    return val


def get_messages(queue: SQSQueue) -> List[SQSMessage]:
    messages = []
    while (message := queue.receive_messages(
        VisibilityTimeout=90,
        WaitTimeSeconds=10,
        MaxNumberOfMessages=1,
    )) != []:
        messages.append(message[0])
    return messages


def handler(event, context):
    rac_bucket = get_env_or_raise("RAC_INPUT_BUCKET")
    project = get_env_or_raise("RAC_PROJECT")
    queue_name = get_env_or_raise("RAC_QUEUE")
    rac_dir = "/tmp"

    sqs = boto3.resource('sqs')
    queue = sqs.get_queue_by_name(QueueName=queue_name)
    messages = get_messages(queue)

    if messages == []:
        raise NoNewFiles(f"Got no new files from {queue_name}")

    s3_client = boto3.client('s3')
    for mess in messages:
        body = json.loads(mess.body)
        key = body["Records"][0]["s3"]["object"]["key"]
        filename = key.split("/")[-1]
        s3_client.download_file(
            Bucket=rac_bucket,
            Key=key,
            Filename=f"{rac_dir}/{filename}",
        )

    subprocess.call([
        "./rac", "-aws", "-project", project, f"{rac_dir}/*.rac",
    ])

    for mess in messages:
        mess.delete()


if __name__ == "__main__":
    handler(None, None)
