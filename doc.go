/*

bounce is a small utility to assist with building HTTP endpoints


PROJECT HOME

See our GitHub repo (https://github.com/atc0005/bounce) for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

PURPOSE

bridge is primarily intended to be used as a HTTP endpoint for testing webhook
payloads. Over time, it may grow other related features to aid in testing
other tools that submit data via HTTP requests.

FEATURES

• single binary, no outside dependencies

• minimal configuration

• Optional submission of client request details to a user-specified Microsoft
  Teams channel (by providing a webhook URL)

• index page automatically generated listing currently supported routes

• request body and associated metadata is echoed to stdout and back to client
  as unformatted request body and automatic formatting of JSON payloads when
  sent to the /api/v1/echo/json endpoint

• Optional, colorization and custom ident control for formatted JSON output

• User configurable TCP port to listen on for incoming HTTP requests

• User configurable IP Address to listen on for incoming HTTP requests

• User configurable logging levels

• User configurable logging format

• User configurable logging output (stdout or stderr)

• User configurable message submission retry and retry delay limits

USAGE

Help output is below. See the README for examples.

$ ./bounce.exe -h

    bounce x.y.z
    https://github.com/atc0005/bounce

    Usage of "T:\github\bounce\bounce.exe":
    -color
            Whether JSON output should be colorized.
    -indent-lvl int
            Number of spaces to use when indenting colorized JSON output. Has no effect unless colorized JSON mode is enabled. (default 2)
    -ipaddr string
            Local IP Address that this application should listen on for incoming HTTP requests. (default "localhost")
    -log-fmt string
            Log messages are written in this format (default "text")
    -log-lvl string
            Log message priority filter. Log messages with a lower level are ignored. (default "info")
    -log-out string
            Log messages are written to this output target (default "stdout")
    -port int
            TCP port that this application should listen on for incoming HTTP requests. (default 8000)
    -webhook-url string
            The Webhook URL provided by a preconfigured Connector. If specified, this application will attempt to send client request details to the Microsoft Teams channel associated with the webhook URL.

*/
package main
