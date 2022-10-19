#!/usr/bin/env python3
import os
import requests
import subprocess
import tarfile
from io import BytesIO

import aws_cdk as cdk

from raclambda.raclambda_stack import RacLambdaStack


RAC_VERSION = "v0.2.6"
RAC_OS = "Linux"
RAC_URL = f"https://github.com/innosat-mats/rac-extract-payload/releases/download/{RAC_VERSION}/Rac_for_{RAC_OS}.tar.gz"  # noqa: E501
RAC_DIR = "./handler"
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
RacLambdaStack(app, "RacLambdaStack",
    # If you don't specify 'env', this stack will be environment-agnostic.
    # Account/Region-dependent features and context lookups will not work,
    # but a single synthesized template can be deployed anywhere.

    # Uncomment the next line to specialize this stack for the AWS Account
    # and Region that are implied by the current CLI configuration.

    #env=cdk.Environment(account=os.getenv('CDK_DEFAULT_ACCOUNT'), region=os.getenv('CDK_DEFAULT_REGION')),

    # Uncomment the next line if you know exactly what Account and Region you
    # want to deploy the stack to. */

    #env=cdk.Environment(account='123456789012', region='us-east-1'),

    # For more information, see https://docs.aws.amazon.com/cdk/latest/guide/environments.html
)

app.synth()
