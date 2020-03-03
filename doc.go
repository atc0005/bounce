/*

bounce is a small utility to assist with building HTTP endpoints


PROJECT HOME

See our GitHub repo (https://github.com/atc0005/bounce) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

PURPOSE

bridge is primarily intended to be used as a HTTP endpoint for testing webhook
payloads. Over time, it may grow other related features to aid in testing
other tools that submit data via HTTP requests.

FEATURES

• single binary, no outside dependencies

• minimal configuration

• index page automatically generated listing currently supported routes

• request body and associated metadata is echoed to stdout and back to client
  • unformatted request body
  • automatic formatting of JSON payloads when sent to the /api/v1/echo/json
	endpoint

• User configurable TCP port to listen on for incoming HTTP requests

• User configurable IP Address to listen on for incoming HTTP requests

USAGE

Help output is below. See the README for examples.

$ ./bounce.exe -h

	2020/03/03 06:28:11 DEBUG: Initializing application

	bounce x.y.z
	https://github.com/atc0005/bounce

	Usage of "T:\github\bounce\bounce.exe":
	-ipaddr string
			Local IP Address that this application should listen on for incoming HTTP requests. (default "localhost")
	-port int
			TCP port that this application should listen on for incoming HTTP requests. (default 8000)

*/
package main
