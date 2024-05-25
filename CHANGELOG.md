# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/) and this project adheres to [Semantic Versioning](http://semver.org/).

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
