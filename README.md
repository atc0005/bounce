# bounce

Small utility to assist with building HTTP endpoints

[![Latest Release](https://img.shields.io/github/release/atc0005/bounce.svg?style=flat-square)](https://github.com/atc0005/bounce/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/bounce?status.svg)](https://godoc.org/github.com/atc0005/bounce)
![Validate Codebase](https://github.com/atc0005/bounce/workflows/Validate%20Codebase/badge.svg)
![Validate Docs](https://github.com/atc0005/bounce/workflows/Validate%20Docs/badge.svg)

- [bounce](#bounce)
  - [Project home](#project-home)
  - [Overview](#overview)
  - [Status](#status)
  - [Features](#features)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
  - [How to install it](#how-to-install-it)
  - [How to use it](#how-to-use-it)
  - [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/bounce) for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

## Overview

This application is primarily intended to be used as a HTTP endpoint for
testing webhook payloads. Over time, it may grow other related features to aid
in testing other tools that submit data via HTTP requests.

## Status

**Under development.**

While usable, the edges are rough and behavior is subject to change.

## Features

- User configurable port to listen on for incoming HTTP requests
- Default "home" or "frontpage" for this application is rendered from either
  default `README.md` file in this repo or user-specified Markdown file
  - Note: Sanitization of Markdown content is applied by default, but this can
    be disabled by command-line flag if desired

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

- Go 1.13+ (for building)
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

Tested using:

- Go 1.13+
- Windows 10 Version 1903
  - native
  - WSL
- Ubuntu Linux 18.04

## How to install it

1. [Download](https://golang.org/dl/) Go
1. [Install](https://golang.org/doc/install) Go
1. Clone the repo
   1. `cd /tmp` (or equivalent)
   1. `git clone https://github.com/atc0005/bounce`
   1. `cd bounce`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     - `sudo yum install make gcc`
   - for Windows
     - Easiest options
       - Skip all of this and build using the default `go build` command in
         Windows
       - build using WSL Ubuntu environment and just copy out the Windows
         binaries from that environment
     - see the StackOverflow Question `32127524` link in the
       [References](#references) section for potential options for installing
       `make` on Windows
     - see the mingw-w64 project homepage link in the
       [References](#references) section for options for installing `gcc` and
       related packages on Windows
1. Build an executable ...
   - for the current operating system (with default `go` build options)
     - `go build`
   - for *all* supported platforms
      - `make all`
   - for Windows only
      - `make windows`
   - for Linux only
     - `make linux`
1. Copy the newly compiled binary to whatever systems that need to run it
   1. Linux: `/tmp/bounce/bounce`
   1. Windows: `/tmp/bounce/bounce.exe`

## How to use it

1. Build via `go build`
1. Open `8080/tcp` in your local firewall to the remote sender
   - Note: Skip this step if you plan to only submit HTTP requests to this
     application running on `localhost`

## References

- `make` on Windows
  - <https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows>
- `gcc` on Windows
  - <https://en.wikipedia.org/wiki/MinGW>
  - <http://mingw-w64.org/>
  - <https://www.msys2.org/>
