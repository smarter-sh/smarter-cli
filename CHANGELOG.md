# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).


## [0.2.1](https://github.com/smarter-sh/smarter-cli/compare/v0.2.0...v0.2.1) (2026-02-04)


### Bug Fixes

* set localhost port to 9357 ([a480b5d](https://github.com/smarter-sh/smarter-cli/commit/a480b5d4384a3591f7bd9a3c12cd254a6b11cbda))

## v0.1.2 (9-Feb-2025)

- fix: change session_id to session_key [c2145f](https://github.com/smarter-sh/smarter-cli/commit/7942589cf86c20f1997f4405a7066e97dec2145f)
- fix: send entire json response to console output formatters [7091f3](https://github.com/smarter-sh/smarter-cli/commit/8890bff142fc202389f45c0684c1da53d27091f3)
- chore: add default values to config on new installs [65b310](https://github.com/smarter-sh/smarter-cli/commit/a2c12caba3802b8a440c443c6c1667dc5165b310)
- fix: incorrect viper references to environment config setting [5af76a](https://github.com/smarter-sh/smarter-cli/commit/f91a595ddb699372ac8884eb66cc1fba495af76a)


## v0.0.2 MVP (24-May-2024)

- connectivity to any of local, alpha, beta, next, prod (default)
- local configuration file with apikey and preferences
- outout in json, yaml and tabular formats
- complete yaml manifest support for Account, ApiConnection, ApiKey, ChatBot, ChatHistory, ChatPluginUsage, Chat, ChatToolCalls, Plugins, SqlConnection, Users.
- manifest verbs supported: apply, delete, deploy, describe, get, logs, example manifest, undeploy
- additional commands: configure, status, version, whoami
- automated ci/cd test, build, release
- automated build for Windows, macOS, Linux
- automated package releases to Homebrew, Chocolatey, apt

## v0.0.1 Initial Release (5-May-2024)

Partial feature set used for scaffolding ci/cd workflows and to setup package manager releases in Homebrew, Chocolatey and apt.
