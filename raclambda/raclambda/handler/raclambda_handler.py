import json
import os
import subprocess
from glob import glob as glob
from pathlib import Path
from tempfile import TemporaryDirectory
from typing import Any, Dict, List, Tuple

import boto3

S3Client = Any
Event = Dict[str, Any]
Context = Any


def get_env_or_raise(envvar: str) -> str:
    if (val := os.environ.get(envvar)) is None:
        raise ValueError(f"'{envvar}' not found in env")
    return val


def parse_event_message(event: Event) -> Tuple[List[str], str]:
    message: Dict[str, Any] = json.loads(event["Records"][0]["body"])
    bucket = message["bucket"]
    objects = message["objects"]
    return objects, bucket


def get_all_files(
    s3_client: S3Client,
    bucket_name: str,
    prefix: str = "",
) -> List[str]:
    file_names = []

    default_kwargs = {
        "Bucket": bucket_name,
        "Prefix": prefix
    }
    next_token = ""

    while next_token is not None:
        updated_kwargs = default_kwargs.copy()
        if next_token != "":
            updated_kwargs["ContinuationToken"] = next_token

        response = s3_client.list_objects_v2(**default_kwargs)
        contents = response.get("Contents")

        for result in contents:
            key = result.get("Key")
            file_names.append(key)

        next_token = response.get("NextContinuationToken")

    return file_names


def download_files(
    s3_client: S3Client,
    bucket_name: str,
    path_name: str,
    file_names: List[str],
) -> None:
    local_path = Path(path_name)

    for file_name in file_names:
        file_path = Path.joinpath(local_path, file_name)
        file_path.parent.mkdir(parents=True, exist_ok=True)
        s3_client.download_file(
            bucket_name,
            file_name,
            str(file_path),
        )


def get_updated_files(
    path_name: str,
    old_file_names: List[str],
) -> List[str]:
    old_files = {Path(f).name for f in old_file_names}
    files = set(glob(f"{path_name}/*"))
    return list(files - old_files)


def upload_files(
    s3_client: S3Client,
    bucket_name: str,
    path_name: str,
    file_names: List[str],
) -> None:
    local_path = Path(path_name)
    for file_name in file_names:
        file_path = Path.joinpath(local_path, file_name)
        s3_client.upload_file(
            bucket_name,
            file_name,
            str(file_path),
        )


def handler(event: Event, context: Context):
    project = get_env_or_raise("RAC_PROJECT")
    dregs_bucket = get_env_or_raise("RAC_DREGS")

    with TemporaryDirectory(
        "_rac",
        "/tmp/",
    ) as rac_dir, TemporaryDirectory(
        "_dregs",
        "/tmp/",
    ) as dregs_dir:
        s3_client = boto3.client('s3')

        # Download RAC files
        objects, rac_bucket = parse_event_message(event)
        download_files(s3_client, rac_bucket, rac_dir, objects)

        # Download Dregs
        dregs = get_all_files(s3_client, dregs_bucket)
        download_files(s3_client, dregs_bucket, dregs_dir, dregs)

        # Process RAC files
        subprocess.call([
            "./rac",
            "-aws",
            "-project", project,
            "-dregs", dregs_dir,
            f"{rac_dir}/*.rac",
        ])

        # Upload new Dregs
        new_dregs = get_updated_files(dregs_dir, dregs)
        upload_files(s3_client, dregs_bucket, dregs_dir, new_dregs)


if __name__ == "__main__":
    handler({}, None)
