# mrename

[![test status](https://github.com/d-ashesss/mrename/workflows/test/badge.svg?branch=main)](https://github.com/d-ashesss/mrename/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/d-ashesss/mrename)](https://goreportcard.com/report/github.com/d-ashesss/mrename)
[![MIT license](https://img.shields.io/github/license/d-ashesss/mrename?color=blue)](https://opensource.org/licenses/MIT)
[![MIT license](https://img.shields.io/github/go-mod/go-version/d-ashesss/mrename)](https://github.com/d-ashesss/mrename)
[![MIT license](https://img.shields.io/github/v/tag/d-ashesss/mrename?include_prereleases&sort=semver)](https://github.com/d-ashesss/mrename)
[![MIT license](https://img.shields.io/badge/may%20contain%20cat%20fur-%F0%9F%90%88-blueviolet)](https://github.com/d-ashesss/mrename)


File renaming CLI utility.

Designed to bulk rename files by a pattern computed with any of supported algorithms.

### Usage

```shell script
mrename [options] converter [converter]...
```

### Arguments

- `converter` Which converter or algorithm to use, to rename files, Supported converters:
  - `md5`
  - `sha1`
  - `tolower`
  - `toupper`
  - `jpeg2jpg`

### Options

- `-c`, `--copy` Copy files instead of renaming them
- `-n`, `--dry-run` Do not actually rename the files
- `-v`, `--verbose` Show detailed output
- `-t`, `--target` Specify the target directory
- `-o`, `--output-format` Output renaming results in specified format. Supported formats:
  - `json`

### Warning

There are no default safegurads. Use at your own risk.
