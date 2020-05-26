# RAC Extract Payload

Download releases from:
https://github.com/innosat-mats/rac-extract-payload/releases

# How to use
## Writing to disk

Run the binary (if you are on windows `rac.exe`):

`rac -project test -description some/racs/info.txt some/racs/*.rac`

The `-project` sets output directory in this case.

The `-description` just copies the file into the output directory and renames it to `ABOUT.txt`.

For more information run `rac --help`

## Sending to AWS (avialable from v0.2.0)

### Obtaining and configuring credentials

First you need to register a user following instructions in the invitation. Access is restricted to those involved with M.A.T.S.

After registering a user, you need to generate security keys. The easiest way is probably to first install awscli and then follow these instructions:

https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html#cli-quick-configuration

But you don't need the cli, the same link above describes were the aws credentials need to go to be found, and it is sufficient to create the file manually.

On linux you can install `apt install awscli` for mac follow these instructions https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-mac.html

### Running rac
Run the binary (if you are on windows `rac.exe`):

`rac -aws -project test -description some/racs/info.txt some/racs/*.rac`

The `-project` should be something concise like "binning2002". Avoid using something that start with "test" if it's something that should be kept since we use "test*" while developing and may remove or overwrite such projects.

The `-description` includes a description file to be sent and will appear as `ABOUT.txt`.

For more information run `rac --help`

### Finding the files

Files will be at:

https://s3.console.aws.amazon.com/s3/buckets/mats-l0-artifacts/?region=eu-north-1

The project you gave will be shown as a folder.

Timeseries like data will be in e.g. "HTR.csv"

CCD-metadata will both be in a json per image and for all images in the "CCD.csv"

### Downloading files from AWS

With the aws-cli installed and the credetials file in place, you can easily download files or entire projects on the commandline. The following example downloads the project "MyArchive1_20200511" from S3 to the local directory "test"

```aws s3 cp s3://mats-l0-artifacts/MyArchive1_20200511 ./test --recursive```

Next example downloads all csv-files from the "MyArchive1_20200511" project to you current local directory.

```aws s3 cp s3://mats-lo-artifacts/MyArchive1_20200511 . --recursive --exclude "*" --include "*.csv"```

# Design
[Design map](docs/README.md)
