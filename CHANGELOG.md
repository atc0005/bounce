# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/bounce/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.4.15] - 2021-11-10

### Overview

- Dependency updates
- built using Go 1.16.10
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.8` to `1.16.10`
  - `actions/checkout`
    - `v2.3.4` to `v2.4.0`

## [v0.4.14] - 2021-09-30

### Overview

- Dependency updates
- built using Go 1.16.8
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.7` to `1.16.8`
  - `actions/setup-node`
    - updated from `v2.4.0` to `v2.4.1`

## [v0.4.13] - 2021-08-09

### Overview

- Dependency updates
- built using Go 1.16.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.6` to `1.16.7`
  - `actions/setup-node`
    - updated from `v2.2.0` to `v2.4.0`

## [v0.4.12] - 2021-07-19

### Overview

- Dependency updates
- built using Go 1.16.6
  - **Statically linked**
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- Add "canary" Dockerfile to track stable Go releases, serve as a reminder to
  generate fresh binaries

### Changed

- Dependencies
  - `Go`
    - `1.16.3` to `1.16.6`
  - `atc0005/go-teams-notify`
    - `v2.5.0` to `v2.6.0`
  - `actions/setup-node`
    - `v2.1.5` to `v2.2.0`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

## [v0.4.11] - 2021-04-09

### Overview

- Misc fixes
- Dependency updates
- built using Go 1.16.3

### Changed

- Dependencies
  - Built using Go 1.16.3
    - **Statically linked**
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `actions/setup-node`
    - updated from `v2.1.4` to `v2.1.5`
  - `atc0005/go-teams-notify`
    - updated from `v2.4.2` to `v2.5.0`

### Fixed

- Linting
  - fieldalignment: struct with X pointer bytes could be Y (govet)
  - Replace deprecated linters: maligned, scopelint
  - SA1019: goteamsnotify.IsValidWebhookURL is deprecated: use
    API.ValidateWebhook() method instead. (staticcheck)

## [v0.4.10] - 2021-02-21

### Overview

- Dependency updates
- Bugfixes
- built using Go 1.15.8

### Changed

- Swap out GoDoc badge for pkg.go.dev badge

- dependencies
  - built using Go 1.15.8
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)
  - `atc0005/go-teams-notify`
    - `v2.3.0` to `v2.4.2`
  - `actions/checkout`
    - `v2.3.3` to `v2.3.4`
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`

### Fixed

- Fix explicit exit code handling

## [v0.4.9] - 2020-10-11

### Added

- Binary release
  - Built using Go 1.15.2
  - **Statically linked**
  - Windows
    - x86
    - x64
  - Linux
    - x86
    - x64

### Changed

- Dependencies
  - upgrade `actions/checkout`
    - `v2.3.2` to `v2.3.3`
  - upgrade `actions/setup-node`
    - `v2.1.1` to `v2.1.2`
- Add `-trimpath` build flag
- Restore explicit exit code handling

### Fixed

- Makefile build options do not generate static binaries
- Misc linting errors raised by latest `gocritic` release included with
  `golangci-lint` `v1.31.0`
- Makefile generates checksums with qualified path

## [v0.4.8] - 2020-08-29

### Changed

- Dependencies
  - upgrade `go.mod` Go version
    - `1.13` to `1.14`
  - upgrade `atc0005/go-teams-notify`
    - `v1.3.1-0.20200419155834-55cca556e726` to `v2.3.0`
      - NOTE: This is a significant change reflecting a merge of required
        functionality from the `atc0005/send2teams` project to the
        `atc0005/go-teams-notify` project
  - upgrade `apex/log`
    - `v1.7.0` to `v1.9.0`
  - upgrade `actions/checkout`
    - `v2.3.1` to `v2.3.2`
  - upgrade `atc0005/send2teams`
    - `v0.4.5` to `v0.4.6`
      - since removed

## [v0.4.7] - 2020-08-04

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable" - currently `Go 1.13.14`
      - "stable" - currently `Go 1.14.6`
      - "unstable" - currently `Go 1.15rc1`
    - Makefile is *not* used in this workflow
    - staticcheck linting using latest stable version provided by the
      `atc0005/go-ci` containers

  - Separate Makefile-based linting and building workflow
    - intended to help ensure that local Makefile-based builds that are
      referenced in project README files continue to work as advertised until
      a better local tool can be discovered/explored further
    - use `golang:latest` container to allow for Makefile-based linting
      tooling installation testing since the `atc0005/go-ci` project provides
      containers with those tools already pre-installed
      - linting tasks use container-provided `golangci-lint` config file
        *except* for the Makefile-driven linting task which continues to use
        the repo-provided copy of the `golangci-lint` configuration file

  - Add Quick Validation workflow
    - run on every push, everything else on pull request updates
    - linting via `golangci-lint` only
    - testing
    - no builds

### Changed

- README
  - Link badges to applicable GitHub Actions workflows results

- Linting
  - local
    - `golangci-lint`
      - disable default exclusions
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
        version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

- Dependencies
  - upgrade `actions/setup-node`
    - `v2.1.0` to `v2.1.1`
  - upgrade `actions/setup-go`
    - `v2.1.0` to `v2.1.1`
    - note: since replaced with a Docker container
  - upgrade `apex/log`
    - `v1.6.0` to `v1.7.0`

## [v0.4.6] - 2020-07-19

### Changed

- Dependencies
  - upgrade `atc0005/send2teams`
    - `v0.4.4` to `v0.4.5`
  - upgrade `TylerBrock/colorjson`
    - `v0.0.0-20180527164720-95ec53f28296` to
      `v0.0.0-20200706003622-8a50f05110d2`
  - upgrade `golang/gddo`
    - `v0.0.0-20200324184333-3c2cc9a6329d` to
      `v0.0.0-20200715224205-051695c33a3f`

## [v0.4.5] - 2020-07-19

### Changed

- Dependencies
  - upgrade `apex/log`
    - `v1.4.0` to `v1.6.0`
  - upgrade `actions/setup-go`
    - `v2.0.3` to `v2.1.0`
  - upgrade `actions/checkout`
    - `v2.3.0` to `v2.3.1`
  - upgrade `actions/setup-node`
    - `v2.0.0` to `v2.1.0`

## [v0.4.4] - 2020-06-17

### Changed

- Dependabot
  - Enable GitHub Actions updates

- Update dependencies
  - `apex/log`
    - `v1.3.0` to `v1.4.0`
  - `actions/setup-go`
    - `v1` to `v2.0.3`
  - `actions/checkout`
    - `v1` to `v2.3.0`
  - `actions/setup-node`
    - `v1` to `v2.0.0`

## [v0.4.3] - 2020-06-16

### Changed

- Update dependencies
  - `apex/log`
    - `v1.1.4` to `v1.3.0`
  - `atc0005/send2teams`
    - `v0.4.1` to `v0.4.4`

- enable dependabot updates

### Fixed

- fix typo in project repo URL

## [v0.4.2] - 2020-04-28

### Fixed

- Remove bash shebang from GitHub Actions Workflow files
- Update README to list accurate build/deploy steps based
  on recent restructuring work

### Changes

- Update golangci-lint to v1.25.1
- Remove gofmt and golint as separate checks, enable
  these linters in golangci-lint config

## [v0.4.1] - 2020-04-25

### Changed

- Install specific binary version of golangci-lint instead of building from
  `master`

- Move golangci-lint settings from Makefile to external config file

- Set go modules mode per `go get` command instead of globally when installing
  linting tools

- Using [vendoring](https://golang.org/cmd/go/#hdr-Vendor_Directories)
  - created top-level `vendor` directory using `go mod vendor`
  - updated GitHub Actions Workflow to specify `-mod=vendor` build flag for
    all `go` commands that I know of that respect the flag
  - updated GitHub Actions Workflow to exclude `vendor` directory from
    Markdown file linting to prevent potential linting issues in vendored
    dependencies from affecting our CI checks
  - updated `Makefile` to use `-mod=vendor` where applicable
  - updated `go vet` linting check to use `-mod=vendor`

- Updated dependencies
  - `apex/log`
    - `v1.1.2` to `v1.1.4`
  - `atc0005/send2teams`
    - `v0.4.0` to `v0.4.1`

### Fixed

- Perform `ioutil.ReadAll()` error check immediately instead of after another
  action takes place
  - minor nit, but potential problem in the future

- CHANGELOG
  - fix release section header refs
    - last release didn't include a link to release entry

- Add missing GoDoc coverage for `routes` package

## [v0.4.0] - 2020-04-23

### Added

- Add support for Microsoft Teams notifications
  - configurable retry, retry delay settings
  - rate-limited submissions to help prevent unintentional abuse of remote API
    - currently hard-coded, but will likely expose this as a flag in a future
      release

- Add monitoring/reporting of notification channels with pending items

- Add monitoring/reporting of notification statistics
  - total
  - pending
  - success
  - failure

- Capture `Ctrl+C` and attempt graceful shutdown

- Plumbed `context` throughout majority of application for cancellation and
  timeout functionality
  - still learning proper use of this package, so likely many mistakes that
    will need to be fixed in a future release

- Logging
  - add *many* more debug statements to help with troubleshooting

### Changed

- Dependencies
  - Use `atc0005/go-teams-notify` package
    - fork of original package with current features and some additional
      changes not yet accepted upstream
  - Use `atc0005/send2teams` package
    - provides wrapper for upstream functionality with message retry, delay
      functionality
    - provides formatting helper functions
    - provides additional webhook URL validation
  - Drop indirect dependency
  - Update `golang/gddo`
  - Add commented entries to have Go use local copies of packages for fast
    prototyping work

### Fixed

- GoDoc formatting
  - remove forced line-wrap which resulted in unintentional code block
    formatting of non-code content

- Refactor logging, flag handling
  - not user visible, so not recording as a "change"

- Manually flush `http.ResponseWriter` to (massively) speed up response time
  for client requests

- Move template parsing to `main()` in an effort to speed up endpoint response
  time for client requests

## [v0.3.3] - 2020-03-14

### Fixed

- Fix potential variable shadowing
- Add missing CHANGELOG subsection header

## [v0.3.2] - 2020-03-05

### Fixed

- Fix CHANGELOG sub-bullet format
- Update README to point readers to milestones for current development status
  details

## [v0.3.1] - 2020-03-05

### Fixed

- (GH-14) Fix potential text template variable shadowing
- (GH-15) Extend / Enhance JSON decoding error handling
  - Add `decodeJSONBody()` method and associated `malformedRequest` type
    provided by Alex Edwards (many thanks for sharing!)
    - Article:
      <https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body>
    - License: MIT (same as this codebase)
    - Book: <https://lets-go.alexedwards.net/>
    - Twitter: <https://twitter.com/ajmedwards>
- (GH-16) README updates to cover v0.3.0 changes
  - add new features to summary list
  - remove leveled logging from "TODO" features list

## [v0.3.0] - 2020-03-04

### Added

- (GH-2) Initial implementation of leveled logging using the `apex/log`
  package
  - logging format flag enables matching handler
    - `discard`
    - `text`
    - `cli`
    - `json`
    - `logfmt`
  - logging output flag allows selection between `stdout` and `stderr`
  - logging level flag allows filtering out log messages lower in priority
    than the chosen value
    - `fatal`
    - `error`
    - `warn`
    - `info`
    - `debug`

### Fixed

- (GH-13) Add missing default values for `ReadTimeout` and `WriteTimeout`
  `http.Server` settings

## [v0.2.1] - 2020-03-04

### Fixed

- (GH-11) Prune invalid Go module entries accidentally introduced in prior
  release

## [v0.2.0] - 2020-03-04

### Added

- (GH-9) Optional, formatted, colorized JSON output for JSON-specific endpoint

## [v0.1.1] - 2020-03-03

### Fixed

- Missing date in v0.1.0 release entry for this CHANGELOG file
- (GH-7) GoDoc formatting: collapse unsupported sub-bullets

## [v0.1.0] - 2020-03-03

### Added

Features of the initial prototype:

- Single binary, no outside dependencies

- Minimal configuration
  - User configurable TCP port to listen on for incoming HTTP requests
    (default: `8000`)
  - User configurable IP Address to listen on for incoming HTTP requests
    (default: `localhost`)
  - Index page automatically generates list of currently supported routes with
    detailed descriptions and supported request methods

- Request body and associated metadata is echoed to stdout and back to client
  - Unformatted request body
  - Automatic formatting of JSON payloads when sent to the /api/v1/echo/json
    endpoint

Worth noting (in no particular order):

- Command-line flags support via `flag` standard library package
- Go modules (vs classic `GOPATH` setup)
- GitHub Actions Workflows which apply linting and build checks
- Makefile for general use cases (including local linting)
  - Note: See README for available options if building on Windows

[Unreleased]: https://github.com/atc0005/bounce/compare/v0.4.15...HEAD
[v0.4.15]: https://github.com/atc0005/bounce/releases/tag/v0.4.15
[v0.4.14]: https://github.com/atc0005/bounce/releases/tag/v0.4.14
[v0.4.13]: https://github.com/atc0005/bounce/releases/tag/v0.4.13
[v0.4.12]: https://github.com/atc0005/bounce/releases/tag/v0.4.12
[v0.4.11]: https://github.com/atc0005/bounce/releases/tag/v0.4.11
[v0.4.10]: https://github.com/atc0005/bounce/releases/tag/v0.4.10
[v0.4.9]: https://github.com/atc0005/bounce/releases/tag/v0.4.9
[v0.4.8]: https://github.com/atc0005/bounce/releases/tag/v0.4.8
[v0.4.7]: https://github.com/atc0005/bounce/releases/tag/v0.4.7
[v0.4.6]: https://github.com/atc0005/bounce/releases/tag/v0.4.6
[v0.4.5]: https://github.com/atc0005/bounce/releases/tag/v0.4.5
[v0.4.4]: https://github.com/atc0005/bounce/releases/tag/v0.4.4
[v0.4.3]: https://github.com/atc0005/bounce/releases/tag/v0.4.3
[v0.4.2]: https://github.com/atc0005/bounce/releases/tag/v0.4.2
[v0.4.1]: https://github.com/atc0005/bounce/releases/tag/v0.4.1
[v0.4.0]: https://github.com/atc0005/bounce/releases/tag/v0.4.0
[v0.3.3]: https://github.com/atc0005/bounce/releases/tag/v0.3.3
[v0.3.2]: https://github.com/atc0005/bounce/releases/tag/v0.3.2
[v0.3.1]: https://github.com/atc0005/bounce/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/bounce/releases/tag/v0.3.0
[v0.2.1]: https://github.com/atc0005/bounce/releases/tag/v0.2.1
[v0.2.0]: https://github.com/atc0005/bounce/releases/tag/v0.2.0
[v0.1.1]: https://github.com/atc0005/bounce/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/bounce/releases/tag/v0.1.0
