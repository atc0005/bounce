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

[Unreleased]: https://github.com/atc0005/bounce/compare/v0.4.1...HEAD
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
