import json
import os
import subprocess
from typing import Any, Dict, List, Tuple

import boto3


def get_env_or_raise(envvar: str) -> str:
    if (val := os.environ.get(envvar)) is None:
        raise ValueError(f"'{envvar}' not found in env")
    return val


def parse_event_message(event: Dict[str, Any]) -> Tuple[List[str], str]:
    message: Dict[str, Any] = json.loads(event["Records"][0]["body"])
    bucket = message["bucket"]
    objects = message["objects"]
    return objects, bucket


def handler(event, context):
    project = get_env_or_raise("RAC_PROJECT")
    # TODO: make temporary dir for rac
    rac_dir = "/tmp"
    # TODO: make temporary dir for slask
    slask_dir = "/tmp"

    objects, bucket = parse_event_message(event)

    # TODO: download slask

    s3_client = boto3.client('s3')
    for key in objects:
        filename = key.split("/")[-1]
        s3_client.download_file(
            Bucket=bucket,
            Key=key,
            Filename=f"{rac_dir}/{filename}",
        )

    subprocess.call([
        "./rac",
        "-aws",
        "-project", project,
        "-slask", slask_dir,
        f"{rac_dir}/*.rac",
    ])

    # TODO: upload slask


if __name__ == "__main__":
    handler(None, None)
