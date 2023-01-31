import json
import os
import subprocess
from pathlib import Path
from tempfile import NamedTemporaryFile, TemporaryDirectory
from typing import Any, Dict, List, Tuple

import boto3

BotoClient = Any
S3Client = BotoClient
SSMClient = BotoClient
Event = Dict[str, Any]
Context = Any


def get_env_or_raise(envvar: str) -> str:
    if (val := os.environ.get(envvar)) is None:
        raise ValueError(f"'{envvar}' not found in env")
    return val


def parse_event_message(event: Event) -> Tuple[str, str]:
    message: Dict[str, Any] = json.loads(event["Records"][0]["body"])
    bucket = message["bucket"]
    object = message["object"]
    return object, bucket


def download_file(
    s3_client: S3Client,
    bucket_name: str,
    path_name: str,
    file_name: str,
) -> Path:
    file_path = Path(path_name) / file_name
    file_path.parent.mkdir(parents=True, exist_ok=True)
    s3_client.download_file(
        bucket_name,
        file_name,
        str(file_path),
    )
    return file_path


def get_rclone_config_path(
    ssm_client: SSMClient,
    rclone_config_ssm_name: str
) -> str:
    rclone_config = ssm_client.get_parameter(
        Name=rclone_config_ssm_name, WithDecryption=True
    )["Parameter"]["Value"]

    f = NamedTemporaryFile(buffering=0, delete=False)
    f.write(rclone_config.encode())

    return f.name


def format_rclone_command(
    config_path: str,
    source: str,
    destination: str,
) -> List[str]:
    cmd = [
        "rclone",
        "--config",
        config_path,
        "copy",
        source,
        destination,
        "--size-only"
    ]

    return cmd


def handler(event: Event, context: Context):
    dregs_bucket = get_env_or_raise("RAC_DREGS")
    output_bucket = get_env_or_raise("RAC_OUTPUT")

    with TemporaryDirectory(
        "_rac",
        "/tmp/",
    ) as rac_dir, TemporaryDirectory(
        "_dregs",
        "/tmp/",
    ) as dregs_dir, TemporaryDirectory(
        "_parquet",
        "/tmp/",
    ) as parquet_dir:
        s3_client = boto3.client('s3')

        # Download RAC file
        object, rac_bucket = parse_event_message(event)
        file = download_file(s3_client, rac_bucket, rac_dir, object)

        # Setup rclone
        ssm_client = boto3.client("ssm")
        rclone_config_path = get_rclone_config_path(
            ssm_client,
            get_env_or_raise("RCLONE_CONFIG_SSM_NAME")
        )

        # Download Dregs
        subprocess.call(format_rclone_command(
            rclone_config_path,
            f"S3:{dregs_bucket}",
            dregs_dir,
            sloppy=True,
        ))

        subprocess.call([
            str(Path(__file__).parent / "rac"),
            "-parquet",
            "-project", parquet_dir,
            "-dregs", dregs_dir,
            str(file),
        ])

        # Upload Parquet files
        subprocess.call(format_rclone_command(
            rclone_config_path,
            parquet_dir,
            f"S3:{output_bucket}",
        ))

        # Sync Dregs
        subprocess.call(format_rclone_command(
            rclone_config_path,
            dregs_dir,
            f"S3:{dregs_bucket}",
            sloppy=True,
        ))


if __name__ == "__main__":
    handler({}, None)
