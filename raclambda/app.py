#!/usr/bin/env python3
import os
import requests
import subprocess
import tarfile
from io import BytesIO

import aws_cdk as cdk

from raclambda.raclambda_stack import RacLambdaStack


RAC_VERSION = "v1.3.0"
RAC_OS = "Linux"
RAC_URL = f"https://github.com/innosat-mats/rac-extract-payload/releases/download/{RAC_VERSION}/Rac_for_{RAC_OS}.tar.gz"  # noqa: E501
RAC_DIR = "./raclambda/handler"
RAC_BIN = f"{RAC_DIR}/rac"

if os.path.exists(RAC_BIN) and RAC_VERSION in subprocess.check_output(
    [RAC_BIN, "-version"]
).decode():
    print("rac binary already up to date")
else:
    print("fetching new rac binary")
    resp = requests.get(RAC_URL)

    if resp.status_code != 200:
        raise RuntimeError(
            f"Got bad response {resp.status_code} when fetching binary"
        )

    with tarfile.open(None, "r:gz", BytesIO(resp.content)) as tf:
        tf.extractall(RAC_DIR)

app = cdk.App()
RacLambdaStack(
    app,
    "RacLambdaStack",
    input_bucket_name="ops-payload-level0-source",
    output_bucket_name="ops-payload-level0-v0.2",
    queue_arn_export_name="L0RACFetcherStackOutputQueue",
    config_ssm_name="/rclone/l0-fetcher",
    rclone_arn="arn:aws:lambda:eu-north-1:671150066425:layer:rclone-amd64:1",
)

app.synth()
