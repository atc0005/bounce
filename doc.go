/*
bounce is a small utility to assist with building HTTP endpoints

# Project Home

See our GitHub repo (https://github.com/atc0005/bounce) for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

# Purpose

bounce is primarily intended to be used as a HTTP endpoint for testing webhook
payloads. Over time, it may grow other related features to aid in testing
other tools that submit data via HTTP requests.

# Features

  - single binary, no outside dependencies
  - minimal configuration
  - Optional submission of client request details to a user-specified
    Microsoft Teams channel (by providing a webhook URL)
  - index page automatically generated listing currently supported routes
  - request body and associated metadata are echoed to stdout and back to
    client
  - echoed request details are provided as-is/unformatted when sent to the
    /api/v1/echo/json endpoint
  - JSON payloads to the /api/v1/echo/json endpoint are automatically
    formatted
  - Optional, colorization and custom ident control for formatted JSON output
  - User configurable TCP port to listen on for incoming HTTP requests
  - User configurable IP Address to listen on for incoming HTTP requests
  - User configurable logging levels
  - User configurable logging format
  - User configurable logging output (stdout or stderr)
  - User configurable message submission retry and retry delay limits

# Usage

See our main README for supported settings and examples.
*/
package main
