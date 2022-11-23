# RAC Extract Payload

Download releases from:
https://github.com/innosat-mats/rac-extract-payload/releases

# How to use
## Writing to disk

Run the binary (if you are on windows `rac.exe`):

`rac -project test -description some/racs/info.txt some/racs/*.rac`

The `-project` sets output directory in this case.

The `-stdout` print output instead of writing to disk, ignoring images.

The `-parquet` save converted data in _Parquet_ format rather than _CSV_, _PNG_ and _JSON_.

The `-dregs` option specifies a directory to use for temporary files written when an unfinished multi-packet is found, in order to continue processing it later.

For more information run `rac --help`

# Design
[Design map](docs/README.md)
