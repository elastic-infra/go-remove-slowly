# go-remove-slowly

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/elastic-infra/go-remove-slowly)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/elastic-infra/go-remove-slowly)
![CircleCI](https://img.shields.io/circleci/build/github/elastic-infra/go-remove-slowly)

A tool to remove file slowly, truncating.

Removing large files with normal `rm` command may cause high I/O load.
This tool truncates the target file gradually, and finally removes it to avoid high I/O load.
