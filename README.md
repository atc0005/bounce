# bounce

Small utility to assist with building HTTP endpoints

## Overview

This application is primarily intended to be used as a HTTP endpoint for
testing webhook payloads. Over time, it may grow other related features to aid
in testing other tools that submit data via HTTP requests.

## Status

**Under development.**

While usable, the edges are rough and behavior is subject to change.

## Directions

1. Build via `go build`
1. Open `8080/tcp` in your local firewall to the remote sender
   - Note: Skip this step if you plan to only submit HTTP requests to this
     application running on `localhost`

## References
