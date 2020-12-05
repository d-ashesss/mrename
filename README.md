# mrename

File renaming CLI utility.

Designed to bulk rename files by a pattern computed with any of supported algorithms.

### Usage

```shell script
mrename [options] converter
```

### Arguments

- `converter` Which converter or algorithm to use, to rename files, Supported converters:
  - `md5`

### Options

- `-n`, `--dry-run` Do not actually rename the files
- `-v`, `--verbose` Show detailed output
- `-t`, `--target` Specify the target directory

### Warning

There are no default safegurads. Use at your own risk.
