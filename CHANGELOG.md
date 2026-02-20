# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe/compare/v1.0.0...v2.0.0) (2026-02-20)

### ⚠ BREAKING CHANGES

* Complete module restructure from Snowflake warehouse to pipe

- Replace snowflake_warehouse resource with snowflake_pipe resource
- Convert to single-module repository layout (removed modules/ directory)
- Update Snowflake provider to snowflakedb/snowflake >= 1.0.0
- Add pipe-specific configuration: database, schema, copy_statement, auto_ingest
- Support AWS SNS topic, error integration, and storage integration
- Update examples: single-pipe and multiple-pipes
- Update Terratest integration tests for pipe resources
- Add header comments to all Terraform configuration files

### Features

* convert warehouse module to pipe module with single-module layout ([713bd48](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe/commit/713bd4831ab2be86035829919861b4881aeea04b))

### Bug Fixes

* enable snowflake_pipe preview feature in examples ([5ba63cd](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe/commit/5ba63cd96745accfa196abb93c6e55f5f7b4812a))

## [unreleased]

### 🚀 Features

- [**breaking**] Convert warehouse module to pipe module with single-module layout

### 🐛 Bug Fixes

- Enable snowflake_pipe preview feature in examples

### 📚 Documentation

- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update Snowflake provider version and pipe configuration
## [1.0.0] - 2026-02-09

### 🚀 Features

- Initial release of snowflake-pipe terraform module

### 🐛 Bug Fixes

- Update tests to use external stages with S3 bucket
- *(snowflake-pipe)* Remove extra blank line in main.tf
- Add snowflake_pipe preview feature configuration
- *(snowflake-pipe)* Remove extra blank line in main.tf

### 📚 Documentation

- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]
- Update CHANGELOG.md [skip ci]

### ⚙️ Miscellaneous Tasks

- *(release)* Version 1.0.0 [skip ci]
