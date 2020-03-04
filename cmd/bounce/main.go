// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"errors"
	"flag"
	"fmt"

	//"log"
	"net/http"
	"os"
	"time"

	"github.com/atc0005/bounce/config"
	"github.com/atc0005/bounce/routes"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/multi"
)

// see templates.go for the hard-coded HTML/CSS template used for the index
// page

func main() {

	// 	apex/log Handlers
	// ---------------------------------------------------------
	// cli - human-friendly CLI output
	// discard - discards all logs
	// es - Elasticsearch handler
	// graylog - Graylog handler
	// json - JSON output handler
	// kinesis - AWS Kinesis handler
	// level - level filter handler
	// logfmt - logfmt plain-text formatter
	// memory - in-memory handler for tests
	// multi - fan-out to multiple handlers
	// papertrail - Papertrail handler
	// text - human-friendly colored output
	// delta - outputs the delta between log calls and spinner

	// create error logger that directs messages to stderr
	// create stdout logger specifically for providing status updates (pretty)
	// create logfmt stdout logger for systemd consumption if `-daemon` flag is used

	// seems apex/log is different. We don't seem to be able to create
	// multiple logger objects, each to a different output target?
	log.SetHandler(multi.New(
		logfmt.New(os.Stderr),
		cli.New(os.Stdout),
	))
	log.SetLevel(log.DebugLevel)
	// log.SetHandler(logfmt.New(os.Stderr))

	//log.NewEntry(logfmt.New(os.Stdout))

	// stderrLogger := &log.Logger{}

	// stderrLogger.SetHandler()

	// ctx := log.WithFields(log.Fields{
	// 	"file": "something.png",
	// 	"type": "image/png",
	// 	"user": "tobi",
	// })

	// ctx.Info("Hello")

	log.Debug("Initializing application")

	appConfig, err := config.NewConfig()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("Failed to initialize application: %s", err)
	}

	log.Debugf("AppConfig: %+v\n", appConfig)

	// SETUP ROUTES
	// See handlers.go for handler definitions

	var ourRoutes routes.Routes
	ourRoutes.Add(routes.Route{
		Name:           "index",
		Description:    "Main page, fallback for unspecified routes",
		Pattern:        "/",
		AllowedMethods: []string{http.MethodGet},
		HandlerFunc:    handleIndex(handleIndexTemplate, &ourRoutes),
	})

	ourRoutes.Add(routes.Route{
		Name:           "echo",
		Description:    "Prints received values as-is to stdout and returns them via HTTP response",
		Pattern:        apiV1EchoEndpointPattern,
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
		HandlerFunc:    echoHandler(echoHandlerTemplate, appConfig.ColorizedJSON, appConfig.ColorizedJSONIndent),
	})

	ourRoutes.Add(routes.Route{
		Name:           "echo-json",
		Description:    "Prints formatted JSON response to stdout and via HTTP response",
		Pattern:        apiV1EchoJSONEndpointPattern,
		AllowedMethods: []string{http.MethodPost},
		HandlerFunc:    echoHandler(echoHandlerTemplate, appConfig.ColorizedJSON, appConfig.ColorizedJSONIndent),
	})

	mux := http.NewServeMux()
	ourRoutes.RegisterWithServeMux(mux)

	// Apply "default" timeout settings provided by Simon Frey; override the
	// default "wait forever" configuration.
	// FIXME: Refine these settings to apply values more appropriate for a
	// small-to-medium on-premise API (e.g., not over a public Internet link
	// where clients are expected to be slow)
	httpServer := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		Handler:           mux,
		Addr:              fmt.Sprintf("%s:%d", appConfig.LocalIPAddress, appConfig.LocalTCPPort),
	}

	// listen on specified port on ALL IP Addresses, block until app is terminated
	log.Infof("Listening on %s port %d ",
		appConfig.LocalIPAddress, appConfig.LocalTCPPort)

	// TODO: This can be handled in a cleaner fashion?
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
