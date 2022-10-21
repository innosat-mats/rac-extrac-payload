from fileinput import filename
import os
import datetime as dt
import subprocess

import boto3


MIN_AGE = 6 * 3600  # 6 hours [s]


class ObjectsTooRecent(Exception):
    pass


def get_env_or_raise(envvar: str) -> str:
    if (val := os.environ.get(envvar)) is None:
        raise ValueError(f"'{envvar}' not found in env")
    return val


def get_objects_list(bucket: str, date: dt.date, client):
    paginator = client.get_paginator("list_objects_v2")

    page_iterator = paginator.paginate(Bucket=bucket)
    objects = list(page_iterator.search(
        f"Contents[?contains(Key, `{date.strftime('%Y%m%d')}`)][]"
    ))

    page_iterator = paginator.paginate(Bucket=bucket)
    day_before = date - dt.timedelta(days=1)
    index = []
    for i, o in enumerate(objects_day_before := list(page_iterator.search(
        f"Contents[?contains(Key, `{day_before.strftime('%Y%m%d')}`)][]"
    ))):
        index.append([o["Key"], i])

    if index != []:
        last_index = index.sort()[-1][1]
        return objects.append(objects_day_before[last_index])
    return objects


def handler(event, context):
    rac_bucket = get_env_or_raise("RAC_INPUT_BUCKET")
    project = get_env_or_raise("RAC_PROJECT")
    rac_dir = "/tmp"

    client = boto3.client('s3')
    now = dt.datetime.utcnow().replace(tzinfo=dt.timezone.utc)
    objects = get_objects_list(rac_bucket, now.date(), client)
    for obj in objects:
        if (now - obj.last_modifed).total_seconds() < MIN_AGE:
            raise ObjectsTooRecent
        
    for obj in objects:
        filename = obj["Key"].split("/")[-1]
        client.dowload_file(
            Bucket=rac_bucket,
            Key=obj["Key"],
            Filename=f"{rac_dir}/{filename}",
        )

    subprocess.call([
        "./rac", "-aws", "-project", project, f"{rac_dir}/*.rac",
    ])


if __name__ == "__main__":
    handler(None, None)
