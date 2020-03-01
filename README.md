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
  - [Available Endpoints](#available-endpoints)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
  - [How to install it](#how-to-install-it)
  - [Configuration Options](#configuration-options)
    - [Configuration file](#configuration-file)
    - [Command-line Arguments](#command-line-arguments)
    - [Worth noting](#worth-noting)
  - [How to use it](#how-to-use-it)
    - [General](#general)
    - [Examples](#examples)
      - [Local: View headers submitted by `GET` request using your browser](#local-view-headers-submitted-by-get-request-using-your-browser)
      - [Local: Submit JSON payload using `curl`, receive unformatted response](#local-submit-json-payload-using-curl-receive-unformatted-response)
      - [Local: Submit JSON payload to JSON-specific endpoint, get formatted response](#local-submit-json-payload-to-json-specific-endpoint-get-formatted-response)
  - [References](#references)
    - [Dependencies](#dependencies)
    - [Instruction / Examples](#instruction--examples)

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

- Available endpoints/routes are limited
- Defined behavior for this application is subject to change

## Features

- User configurable port to listen on for incoming HTTP requests

## Available Endpoints

Below is a static listing of the available endpoints that may be used for
testing with this application. Visiting the `index` should also generate a
dynamic listing of the available endpoints. Please [open an
issue](https://github.com/atc0005/bounce/issues) if you find that there is a
mismatch between these entries and those listed on the application `index`.

| Name        | Pattern             | Description                                                                                          | Allowed Methods                | Supported Request content types  | Expected Response content type |
| ----------- | ------------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------ | -------------------------------- | ------------------------------ |
| `index`     | `/`                 | Main page, fallback for unspecified routes.                                                          | `GET`                          | `text/plain`                     | `text/html`                    |
| `echo`      | `/api/v1/echo`      | Prints received values to stdout and returns them (as-is, no "pretty" formatting) via HTTP response. | `GET`, `POST`                  | `text/plain`, `application/json` | `text/plain`                   |
| `echo-json` | `/api/v1/echo/json` | Prints "pretty printed" JSON request body to stdout and returns via HTTP response.                   | `GET` (limited), `POST` (JSON) | `text/plain`, `application/json` | `text/plain`                   |

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

## Configuration Options

### Configuration file

- TODO: Evaluate whether this would be particularly beneficial or if the CLI
  flags are sufficient for our purposes

### Command-line Arguments

| Option      | Required | Default | Repeat | Possible                                | Description                                                                 |
| ----------- | -------- | ------- | ------ | --------------------------------------- | --------------------------------------------------------------------------- |
| `h`, `help` | No       | `false` | No     | `h`, `help`                             | Show Help text along with the list of supported flags.                      |
| `port`      | No       | `8000`  | No     | *valid whole numbers* (see Description) | TCP port that this application should listen on for incoming HTTP requests. |

### Worth noting

- For best results, limit your choice of TCP port to an unprivileged user
  port between `1024` and `49151`

## How to use it

### General

1. Build or obtain a pre-compiled copy of the executable appropriate for your
   operating system
   - NOTE: As of this writing, CI-enabled automatic builds for new releases is
     not yet available. We hope to add this in the near future.
1. Pick a TCP port where you will have the application listen
   - e.g., `8000`
1. Open your chosen TCP port in the firewall on the system where this
   application will run
   - if possible, limit access to just the remote system submitting HTTP
     requests
   - skip this step if you plan to only submit HTTP requests from your own
     system to this application running *on* your system
     - e.g., `localhost:8000`
1. Visit the index page for this application at the appropriate IP Address and
   the port you specified
   - e.g., `http://localhost:8000/`
1. Chose one of the available routes that meet your requirements

### Examples

#### Local: View headers submitted by `GET` request using your browser

```ShellSession
./bounce.exe -port 8000
2020/03/01 04:58:37 DEBUG: Initializing application
2020/03/01 04:58:37 DEBUG: Valid, non-privileged user port between 1024 and 49151 configured: 8000
2020/03/01 04:58:37 DEBUG: LocalTCPPort: 8000
2020/03/01 04:58:37 DEBUG: Add index to routes ...
2020/03/01 04:58:37 DEBUG: Add echo to routes ...
2020/03/01 04:58:37 DEBUG: Register index with ServeMux ...
2020/03/01 04:58:37 DEBUG: Register echo with ServeMux ...
2020/03/01 04:58:37 Listening on port 8000
DEBUG: echoHandler endpoint hit

HTTP Method used by client: GET
Client IP Address: [::1]:52555

Headers:

  * User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4074.0 Safari/537.36
  * Sec-Fetch-Site: none
  * Sec-Fetch-Mode: navigate
  * Sec-Fetch-User: ?1
  * Sec-Fetch-Dest: document
  * Connection: keep-alive
  * Upgrade-Insecure-Requests: 1
  * Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
  * Accept-Encoding: gzip, deflate, br
  * Accept-Language: en-US,en;q=0.9

```

Items to note:

- Port `8000` is the default, we're just being explicit here.
- I ran the application on Windows 10 Version 1903
- I visited the `/echo` endpoint (`http://localhost:8000/echo`) from Google Chrome Canary
- The same headers shown here are also shown in the browser

#### Local: Submit JSON payload using `curl`, receive unformatted response

```ShellSession
$ curl -X POST -H "Content-Type: application/json" -d @contrib/splunk-test-payload.json http://localhost:8000/api/v1/echo
DEBUG: echoHandler endpoint hit

HTTP Method used by client: POST
Client IP Address: 127.0.0.1:53098

Headers:

  * User-Agent: curl/7.68.0
  * Accept: */*
  * Content-Type: application/json
  * Content-Length: 306
POST request received; reading Body value ...
Body:
{    "result": {        "sourcetype": "mongod",        "count": "8"    },    "sid": "scheduler_admin_search_W2_at_14232356_132",    "results_link": "http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132",    "search_name": null,    "owner": "admin",    "app": "search"}
```

Items to note:

- Output shown above is not wrapped
- `curl` was executed from within a `Git Bash` shell session
- The current working directory was the root of the cloned repo
- We used a sample Splunk Webhook request JSON payload provided in the official docs
  - see "Splunk Enterprise > Alerting Manual > Use a webhook alert action"
- Non-plaintext submissions are *not* "pretty-printed" or formatted in any way

#### Local: Submit JSON payload to JSON-specific endpoint, get formatted response

```ShellSession
$ curl -X POST -H "Content-Type: application/json" -d @contrib/splunk-test-payload.json http://localhost:8000/api/v1/echo/json

TODO
```

Note:

- Output was not modified, but copied as-is from the terminal session
- Output was "pretty-printed" by the application
- `curl` was executed from within a `Git Bash` shell session
- The current working directory was the root of the cloned repo

## References

### Dependencies

- `make` on Windows
  - <https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows>
- `gcc` on Windows
  - <https://en.wikipedia.org/wiki/MinGW>
  - <http://mingw-w64.org/>
  - <https://www.msys2.org/>

### Instruction / Examples

- General
  - <https://gobyexample.com/http-servers>
  - <https://stackoverflow.com/questions/24556001/how-to-range-over-slice-of-structs-instead-of-struct-of-slices>
  - <https://golangcode.com/get-the-request-ip-addr/>
  - <https://github.com/eddturtle/golangcode-site>
  - <https://stackoverflow.com/questions/19038598/how-can-i-pretty-print-json-using-go/42426889>
  - <https://stackoverflow.com/questions/43021058/golang-read-request-body/43021236#43021236>
  - <https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body>

- Splunk / JSON payload
  - [Splunk Enterprise (v8.0.1) > Alerting Manual > Use a webhook alert action](https://docs.splunk.com/Documentation/Splunk/8.0.1/Alert/Webhooks)
