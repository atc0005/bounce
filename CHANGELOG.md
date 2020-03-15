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

[Unreleased]: https://github.com/atc0005/bounce/compare/v0.3.3...HEAD
[v0.3.3]: https://github.com/atc0005/bounce/releases/tag/v0.3.3
[v0.3.2]: https://github.com/atc0005/bounce/releases/tag/v0.3.2
[v0.3.1]: https://github.com/atc0005/bounce/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/bounce/releases/tag/v0.3.0
[v0.2.1]: https://github.com/atc0005/bounce/releases/tag/v0.2.1
[v0.2.0]: https://github.com/atc0005/bounce/releases/tag/v0.2.0
[v0.1.1]: https://github.com/atc0005/bounce/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/bounce/releases/tag/v0.1.0
