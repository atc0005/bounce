# bounce

Small utility to assist with building HTTP endpoints

[![Latest Release](https://img.shields.io/github/release/atc0005/bounce.svg?style=flat-square)][latest-release]
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/bounce.svg)](https://pkg.go.dev/github.com/atc0005/bounce)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/bounce)](https://github.com/atc0005/bounce)
[![Lint and Build](https://github.com/atc0005/bounce/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/bounce/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/bounce/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/bounce/actions/workflows/project-analysis.yml)

- [bounce](#bounce)
  - [Project home](#project-home)
  - [Overview](#overview)
  - [Features](#features)
    - [Current](#current)
    - [Future](#future)
  - [Available Endpoints](#available-endpoints)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
    - [Building source code](#building-source-code)
    - [Running](#running)
  - [Installation](#installation)
    - [From source](#from-source)
    - [Using release binaries](#using-release-binaries)
  - [Configuration Options](#configuration-options)
    - [Configuration file](#configuration-file)
    - [Command-line Arguments](#command-line-arguments)
    - [Worth noting](#worth-noting)
  - [How to use it](#how-to-use-it)
    - [General](#general)
    - [Examples](#examples)
      - [Local: Send client request details to Microsoft Teams](#local-send-client-request-details-to-microsoft-teams)
      - [Local: View headers submitted by `GET` request using your browser](#local-view-headers-submitted-by-get-request-using-your-browser)
      - [Local: Submit JSON payload using `curl`, receive unformatted response](#local-submit-json-payload-using-curl-receive-unformatted-response)
      - [Local: Submit JSON payload using `curl` to JSON-specific endpoint, get formatted response](#local-submit-json-payload-using-curl-to-json-specific-endpoint-get-formatted-response)
      - [Local: Submit JSON payload using `curl` to JSON-specific endpoint, get colorized, formatted response](#local-submit-json-payload-using-curl-to-json-specific-endpoint-get-colorized-formatted-response)
  - [References](#references)
    - [Dependencies](#dependencies)
    - [Instruction / Examples](#instruction--examples)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Overview

This application is primarily intended to be used as a HTTP endpoint for
testing webhook payloads. Over time, it may grow other related features to aid
in testing other tools that submit data via HTTP requests.

## Features

### Current

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
  - Optional, colorization and custom ident control for formatted JSON output

- Optional submission of client request details to a user-specified Microsoft
  Teams channel (by providing a webhook URL)

- User configurable logging settings
  - levels, format and output (see command-line arguments table)

- Message delivery retry support with retry and retry delay values
  configurable via flag
  - currently used by Microsoft Teams notifications support, also intended for
    use with future email notifications support

- Capture `Ctrl+C` and attempt graceful shutdown

- Notification statistics emitted periodically to assist with troubleshooting

### Future

| Priority | Milestone                                                         | Description                                                                                                                     |
| -------- | ----------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| Low      | [Unplanned](https://github.com/atc0005/bounce/milestone/8)        | Potential merit, but are either low demand or are more complex to implement than other issues under consideration.              |
| Medium   | [Future](https://github.com/atc0005/bounce/milestone/2)           | Considered to have merit and sufficiently low complexity to fit within a near-future milestone.                                 |
| High     | [vX.Y.Z](https://github.com/atc0005/bounce/milestones?state=open) | Milestones with a semantic versioning pattern reflect collections of issues that are in a planning or active development state. |

## Available Endpoints

Below is a static listing of the available endpoints that may be used for
testing with this application. Visiting the `index` should also generate a
dynamic listing of the available endpoints. Please [open an
issue](https://github.com/atc0005/bounce/issues) if you find that there is a
mismatch between these entries and those listed on the application `index`.

| Name        | Pattern             | Description                                                                        | Allowed Methods                | Supported Request content types  | Expected Response content type |
| ----------- | ------------------- | ---------------------------------------------------------------------------------- | ------------------------------ | -------------------------------- | ------------------------------ |
| `index`     | `/`                 | Main page, fallback for unspecified routes.                                        | `GET`                          | `text/plain`                     | `text/html`                    |
| `echo`      | `/api/v1/echo`      | Prints received values as-is to stdout and returns them via HTTP response.         | `GET`, `POST`                  | `text/plain`, `application/json` | `text/plain`                   |
| `echo-json` | `/api/v1/echo/json` | Prints "pretty printed" JSON request body to stdout and returns via HTTP response. | `GET` (limited), `POST` (JSON) | `text/plain`, `application/json` | `text/plain`                   |

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for *preferred* version
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 10
- Ubuntu Linux 18.04+

## Installation

### From source

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
   - NOTE: Pay special attention to the remarks about `$HOME/.profile`
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
     - Emulated environments (*easier*)
       - Skip all of this and build using the default `go build` command in
         Windows
       - build using Windows Subsystem for Linux Ubuntu environment and just
         copy out the Windows binaries from that environment
       - If already running a Docker environment, use a container with the Go
         tool-chain already installed
       - If already familiar with LXD, create a container and follow the
         installation steps given previously to install required dependencies
     - Native tooling (*harder*)
       - see the StackOverflow Question `32127524` link in the
         [References](#references) section for potential options for
         installing `make` on Windows
       - see the mingw-w64 project homepage link in the
         [References](#references) section for options for installing `gcc`
         and related packages on Windows
1. Build
   - for current operating system
     - `go build -mod=vendor ./cmd/bounce/`
       - *forces build to use bundled dependencies in top-level `vendor`
         folder*
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Copy the applicable binary to whatever systems needs to run it
   - if using `Makefile`: look in `/tmp/release_assets/bounce/`
   - if using `go build`: look in `/tmp/bounce/`

**NOTE**: Depending on which `Makefile` recipe you use the generated binary
may be compressed and have an `xz` extension. If so, you should decompress the
binary first before deploying it (e.g., `xz -d bounce-linux-amd64.xz`).

### Using release binaries

1. Download the [latest release][latest-release] binaries
1. Decompress binaries
   - e.g., `xz -d bounce-linux-amd64.xz`
1. Deploy
   - Place in a location of your choice

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

## Configuration Options

### Configuration file

- TODO: Evaluate whether this would be particularly beneficial or if the CLI
  flags are sufficient for our purposes

### Command-line Arguments

| Option          | Required | Default        | Repeat | Possible                                   | Description                                                                                                                                                                                       |
| --------------- | -------- | -------------- | ------ | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `h`, `help`     | No       | `false`        | No     | `h`, `help`                                | Show Help text along with the list of supported flags.                                                                                                                                            |
| `port`          | No       | `8000`         | No     | *valid whole numbers*                      | TCP port that this application should listen on for incoming HTTP requests.                                                                                                                       |
| `ipaddr`        | No       | `localhost`    | No     | *valid fqdn, local name or IP Address*     | Local IP Address that this application should listen on for incoming HTTP requests.                                                                                                               |
| `color`         | No       | `false`        | No     | `true`, `false`                            | Whether JSON output should be colorized.                                                                                                                                                          |
| `indent-lvl`    | No       | `2`            | No     | *1+; positive whole numbers*               | Number of spaces to use when indenting colorized JSON output. Has no effect unless colorized JSON mode is enabled.                                                                                |
| `log-lvl`       | No       | `info`         | No     | `fatal`, `error`, `warn`, `info`, `debug`  | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                         |
| `log-out`       | No       | `stdout`       | No     | `stdout`, `stderr`                         | Log messages are written to this output target.                                                                                                                                                   |
| `log-fmt`       | No       | `text`         | No     | `cli`, `json`, `logfmt`, `text`, `discard` | Use the specified `apex/log` package "handler" to output log messages in that handler's format.                                                                                                   |
| `webhook-url`   | No       | *empty string* | No     | *valid webhook URL*                        | The Webhook URL provided by a preconfigured Connector. If specified, this application will attempt to send client request details to the Microsoft Teams channel associated with the webhook URL. |
| `retries`       | No       | `2`            | No     | *positive whole number*                    | The number of attempts that this application will make to deliver messages before giving up.                                                                                                      |
| `retries-delay` | No       | `2`            | No     | *positive whole number*                    | The number of seconds that this application will wait before making another delivery attempt.                                                                                                     |

### Worth noting

- For best results, limit your choice of TCP port to an unprivileged user
  port between `1024` and `49151`

- Log format names map directly to the Handlers provided by the `apex/log`
  package. Their descriptions are duplicated below from the [official
  README](https://github.com/apex/log/blob/master/Readme.md) for reference:

| Log Format ("Handler") | Description                        |
| ---------------------- | ---------------------------------- |
| `cli`                  | human-friendly CLI output          |
| `json`                 | provides log output in JSON format |
| `logfmt`               | plain-text logfmt output           |
| `text`                 | human-friendly colored output      |
| `discard`              | discards all logs                  |

- Microsoft Teams webhook URLs have one of two known prefixes. Both are valid
  as of this writing, but new webhook URLs only appear to be generated using
  the first prefix.
  1. <https://outlook.office.com>
  1. <https://outlook.office365.com>

## How to use it

### General

1. Build or obtain a pre-compiled copy of the executable appropriate for your
   operating system
   - NOTE: As of this writing, CI-enabled automatic builds for new releases
     are not yet available. We hope to add this in the near future.
1. Pick a TCP port where you will have the application listen
   - e.g., `8000`
1. Decide what IP Address that you wish to have this application "bind" or
   "listen" on
   - e.g., `localhost` or `192.168.1.100` (*arbitrary number shown here*)
1. Update your host firewall on the system where this application will run to
   permit connections to your chosen IP Address and TCP port
   - if possible, limit access to just the remote system submitting HTTP
     requests
   - skip this step if you plan to only submit HTTP requests from your own
     system to this application running *on* your system
     - e.g., `localhost:8000`
1. Run this application using your preferred settings by specifying the
   appropriate command-line flag options.
   - e.g., if you specify a valid Outlook/Microsoft Teams webhook URL, this
     application will attempt to send client request details to the associated
     Microsoft Teams channel.
1. Visit the index page for this application at the appropriate IP Address and
   the port you specified
   - e.g., `http://localhost:8000/`
1. Chose one of the available routes that meet your requirements

### Examples

#### Local: Send client request details to Microsoft Teams

In order to have messages sent to a Microsoft Teams channel, `bounce` requires
that you provide a webhook URL to use for submissions. Many more options are
available. See the list of available flags for details.

Short example:

`./bounce.exe -webhook-url WEBHOOK_URL_HERE`

Full example:

```ShellSession
$ ./bounce.exe -webhook-url "https://outlook.office.com/webhook/a1269812-6d10-44b1-abc5-b84f93580ba0@9e7b80c7-d1eb-4b52-8582-76f921e416d9/IncomingWebhook/3fdd6767bae44ac58e5995547d66a4e4/f332c8d9-3397-4ac5-957b-b8e3fc465a8c"

  INFO[0000] bounce is listening on localhost port 8000
  INFO[0000] Visit http://localhost:8000 in your web browser for details
```

When visiting <http://localhost:8000/api/v1/echo> in your browser the
following will be logged in the console:

```log
Request received: 2020-04-23 16:42:27
Endpoint path requested by client: /api/v1/echo
HTTP Method used by client: GET
Client IP Address: 127.0.0.1:53093

Headers:


  * Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
  * Accept-Encoding: gzip, deflate, br
  * Accept-Language: en-US,en;q=0.9
  * Cache-Control: max-age=0
  * Connection: keep-alive
  * Referer: http://localhost:8000/
  * Sec-Fetch-Dest: document
  * Sec-Fetch-Mode: navigate
  * Sec-Fetch-Site: same-origin
  * Sec-Fetch-User: ?1
  * Upgrade-Insecure-Requests: 1
  * User-Agent: Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4124.2 Safari/537.36



No request body was provided by client.





  INFO[0405] StartNotifyMgr: sendMessage: Message successfully sent to Microsoft Teams at 16:42:33
```

and finally this is what I see in a test Microsoft Teams channel:

<!-- Attempt to use image reference and inherit the alt-text already set -->
![Microsoft Teams test submission example screenshot for v0.4.0
release][screenshot-microsoft-teams-example-submission-v0.4.0]

Note:

- Output is from preparing the `v0.4.0` release and is subject to change
- The `go-teams-notify` named webhook connector is a test entry initially
  created when testing the `dasrick/go-teams-notify` package; I've just
  kept using this test connector ever since
- Port `8000` is the default, but can be overridden via a command-line flag
- Localhost is the default, but another active IP Address can be specified via
  a command-line flag
- I ran the application on Windows 10 Version 1903
- I visited the `/echo` endpoint (`http://localhost:8000/echo`) from Google Chrome Canary
- The same non-logging output shown here is also shown in the browser

#### Local: View headers submitted by `GET` request using your browser

```ShellSession
$ ./bounce.exe

  INFO[0000] Listening on localhost port 8000

Request received: 2020-03-04T22:03:31-06:00
Endpoint path requested by client: /api/v1/echo
HTTP Method used by client: GET
Client IP Address: 127.0.0.1:60465

Headers:


  * Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
  * Accept-Encoding: gzip, deflate, br
  * Accept-Language: en-US,en;q=0.9
  * Connection: keep-alive
  * Referer: http://localhost:8000/
  * Sec-Fetch-Dest: document
  * Sec-Fetch-Mode: navigate
  * Sec-Fetch-Site: same-origin
  * Sec-Fetch-User: ?1
  * Upgrade-Insecure-Requests: 1
  * User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/82.0.4077.0 Safari/537.36



No request body was provided by client.

```

Note:

- Output is from a `v0.3.0` release and is subject to change
- Port `8000` is the default, we're just being explicit here.
- I ran the application on Windows 10 Version 1903
- I visited the `/echo` endpoint (`http://localhost:8000/echo`) from Google Chrome Canary
- The same non-logging output shown here is also shown in the browser

#### Local: Submit JSON payload using `curl`, receive unformatted response

Start `bounce` from one terminal:

```ShellSession
$ ./bounce.exe

  INFO[0000] Listening on localhost port 8000
```

Submit a `curl` request from another to a non-JSON specific endpoint and see
the same output from `bounce` or from the terminal where you ran the `curl`
command:

```ShellSession
$ curl --silent -X POST -H "Content-Type: application/json" -d @contrib/splunk-test-payload-unformatted.json http://localhost:8000/api/v1/echo

Request received: 2020-03-04T22:06:48-06:00
Endpoint path requested by client: /api/v1/echo
HTTP Method used by client: POST
Client IP Address: 127.0.0.1:60498

Headers:


  * Accept: */*
  * Content-Length: 254
  * Content-Type: application/json
  * User-Agent: curl/7.68.0



Unformatted request body:

{"result":{"sourcetype":"mongod","count":"8"},"sid":"scheduler_admin_search_W2_at_14232356_132","results_link":"http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132","search_name":null,"owner":"admin","app":"search"}




Error formatting request body:

This endpoint does not apply JSON formatting to the request body.
Use the "/api/v1/echo/json" endpoint for JSON payload testing.

```

Note:

- Output is from a `v0.3.0` release and is subject to change
- Output shown above is not wrapped
- We used a "minified" version of the sample Splunk Webhook request JSON
  payload found in the official docs which is *not* wrapped or formatted
  - see "Splunk Enterprise > Alerting Manual > Use a webhook alert action"
- `curl` was executed from within a `Git Bash` shell session
- The current working directory was the root of the cloned repo
- Non-plaintext submissions are *not* "pretty-printed" or formatted in any way

#### Local: Submit JSON payload using `curl` to JSON-specific endpoint, get formatted response

Start `bounce` from one terminal:

```ShellSession
$ ./bounce.exe

  INFO[0000] Listening on localhost port 8000
```

Submit a `curl` request from another to a JSON-specific endpoint and see
the same output from `bounce` or from the terminal where you ran the `curl`
command:

```ShellSession
$ curl --silent -X POST -H "Content-Type: application/json" -d @contrib/splunk-test-payload-unformatted.json http://localhost:8000/api/v1/echo/json

Request received: 2020-03-04T22:05:47-06:00
Endpoint path requested by client: /api/v1/echo/json
HTTP Method used by client: POST
Client IP Address: 127.0.0.1:60497

Headers:


  * Accept: */*
  * Content-Length: 254
  * Content-Type: application/json
  * User-Agent: curl/7.68.0



Unformatted request body:

{"result":{"sourcetype":"mongod","count":"8"},"sid":"scheduler_admin_search_W2_at_14232356_132","results_link":"http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132","search_name":null,"owner":"admin","app":"search"}



Formatted Body:

{
        "result": {
                "sourcetype": "mongod",
                "count": "8"
        },
        "sid": "scheduler_admin_search_W2_at_14232356_132",
        "results_link": "http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132",
        "search_name": null,
        "owner": "admin",
        "app": "search"
}

```

Note:

- Output is from a `v0.3.0` release and is subject to change
- Output was not modified, but copied as-is from the terminal session
- Output was formatted or "pretty-printed" by the application
- `curl` was executed from within a `Git Bash` shell session
- The current working directory was the root of the cloned repo

#### Local: Submit JSON payload using `curl` to JSON-specific endpoint, get colorized, formatted response

Same as our other JSON-specific endpoint example, but with colorized output enabled.

Here is what you get without color:

<!-- Attempt to use image reference and inherit the alt-text already set -->
![Uncolored JSON output example screenshot][screenshot-uncolored-json-output]

and with colorized JSON output enabled:

<!-- Attempt to use image reference and inherit the alt-text already set -->
![Colored JSON output example screenshot for v0.2.0 release][screenshot-colored-json-output-v0.2.0]

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
  - <https://stackoverflow.com/questions/22886598/how-to-handle-errors-in-goroutines>
    - <https://stackoverflow.com/a/22887491>
  - <https://groups.google.com/forum/#!topic/golang-nuts/QEORIGKZO24>
    - explains benefits of 1-deep buffered channels (asynchronous) vs
      unbuffered (synchronous)
      - Bakul Shah: *In general, synchronize only when you have to. Here the main thread
        wants to know when the worker thread terminates but the worker thread
        doesn't care when the main thread gets around to reading from "done".
        Using a 1 deep buffer channel exactly captures this usage pattern.  An
        unbuffered channel would make the worker thread "rendezvous" with the
        main thread, which is unnecessary.*
  - <https://golang.org/ref/spec#Length_and_capacity>
  - <https://gobyexample.com/closures>
  - <https://golangr.com/closure/>

- Contexts
  - <https://gobyexample.com/context>
  - <https://gobyexample.com/timeouts>
  - <https://golang.org/pkg/context/#WithCancel>
  - <https://groups.google.com/forum/#!topic/golang-nuts/IJXjldvpNQM>
  - <https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97>
  - <https://marcofranssen.nl/go-webserver-with-graceful-shutdown/>

- Request body
  - <https://stackoverflow.com/questions/43021058/golang-read-request-body/43021236#43021236>
  - <https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body>
  - <https://stackoverflow.com/questions/33532374/in-go-how-can-i-reuse-a-readcloser>

- HTTP Server
  - <https://blog.simon-frey.eu/go-as-in-golang-standard-net-http-config-will-break-your-production/>
  - <https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779>
  - <https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97>
  - <https://marcofranssen.nl/go-webserver-with-graceful-shutdown/>

- Logging
  - <https://github.com/apex/log>
  - <https://brandur.org/logfmt>

- Formatted / Colored JSON
  - <https://stackoverflow.com/questions/19038598/how-can-i-pretty-print-json-using-go/42426889>
  - <https://stackoverflow.com/a/50549770/903870>
  - <https://github.com/TylerBrock/colorjson>

- Splunk / JSON payload
  - [Splunk Enterprise (v8.0.1) > Alerting Manual > Use a webhook alert action](https://docs.splunk.com/Documentation/Splunk/8.0.1/Alert/Webhooks)

- Microsoft Teams
  - <https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using#rate-limiting-for-connectors>

<!-- Screenshot references for use within example section  -->
[screenshot-uncolored-json-output]: media/v0.2.0/bounce-json-uncolored-output-2020-03-04.png "Uncolored JSON output example screenshot"
[screenshot-colored-json-output-v0.2.0]: media/v0.2.0/bounce-json-colorizer-output-2020-03-04.png "Colored JSON output example screenshot for v0.2.0 release"
[screenshot-microsoft-teams-example-submission-v0.4.0]: media/v0.4.0/microsoft-teams-test-submission-2020-04-23.png "Microsoft Teams test submission example screenshot for v0.4.0 release"

<!-- Additional footnotes here  -->

[repo-url]: <https://github.com/atc0005/bounce>  "This project's GitHub repo"

[latest-release]: <https://github.com/atc0005/bounce/releases/latest> "Latest release for this project"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
