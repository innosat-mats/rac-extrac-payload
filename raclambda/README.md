# RAC Lambda
RAC Lamda is a lambda function that regularly checks for file updates from a
specified bucket. The bucket should be the one used by
[L0-fetcher][https://github.com/innosat-mats/L0-fetcher]. Files are downloaded
locally to the lambda and processed by the `rac` binary. The delivery
destination is decided by the `rac` binary. The stack sets up lambda permissions
for the bucket used for outputting artifacts so it is important that these are
the same in both the stack and the binary.

## RAC binary dependency
The lambda requires an up to date `rac` binary. The synth-step will
automatically download a rac binary of the version specified in `app.py`.

## Deploy
1. Make sure your aws credentials are set up properly

2. Make sure deployment parameters in `app.py` and arguments to the
   `RacLambdaStack` are what you want.

3. Run `cdk deploy`

## Test
Tests are run by running `tox`.

## Useful commands

 * `cdk ls`          list all stacks in the app
 * `cdk synth`       emits the synthesized CloudFormation template
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk docs`        open CDK documentation

Enjoy!
