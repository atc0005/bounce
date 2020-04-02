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
	"net/http"
	"os"

	"github.com/atc0005/bounce/config"
	"github.com/atc0005/bounce/routes"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"
)

// for handler in cli discard es graylog json kinesis level logfmt memory
// multi papertrail text delta; do go get
// github.com/apex/log/handlers/${handler}; done

// see templates.go for the hard-coded HTML/CSS template used for the index
// page

func main() {

	// This will use default logging settings (level filter, destination)
	// as the application hasn't "booted up" far enough to apply custom
	// choices yet.
	log.Debug("Initializing application")

	appConfig, err := config.NewConfig()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("Failed to initialize application: %s", err)
	}

	var logOutput *os.File
	switch appConfig.LogOutput {
	case config.LogOutputStderr:
		logOutput = os.Stderr
	case config.LogOutputStdout:
		logOutput = os.Stdout
	}

	switch appConfig.LogFormat {
	case config.LogFormatCLI:
		log.SetHandler(cli.New(logOutput))
	case config.LogFormatJSON:
		log.SetHandler(json.New(logOutput))
	case config.LogFormatLogFmt:
		log.SetHandler(logfmt.New(logOutput))
	case config.LogFormatText:
		log.SetHandler(text.New(logOutput))
	case config.LogFormatDiscard:
		log.SetHandler(discard.New())
	}

	switch appConfig.LogLevel {
	case config.LogLevelFatal:
		log.SetLevel(log.FatalLevel)
	case config.LogLevelError:
		log.SetLevel(log.ErrorLevel)
	case config.LogLevelWarn:
		log.SetLevel(log.WarnLevel)
	case config.LogLevelInfo:
		log.SetLevel(log.InfoLevel)
	case config.LogLevelDebug:
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("AppConfig: %+v", appConfig)

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
		HandlerFunc: echoHandler(
			echoHandlerTemplate,
			appConfig.ColorizedJSON,
			appConfig.ColorizedJSONIndent,
			appConfig.WebhookURL,
		),
	})

	mux := http.NewServeMux()
	ourRoutes.RegisterWithServeMux(mux)

	// Apply "default" timeout settings provided by Simon Frey; override the
	// default "wait forever" configuration.
	// FIXME: Refine these settings to apply values more appropriate for a
	// small-to-medium on-premise API (e.g., not over a public Internet link
	// where clients are expected to be slow)
	httpServer := &http.Server{
		ReadHeaderTimeout: config.HTTPServerReadHeaderTimeout,
		ReadTimeout:       config.HTTPServerReadTimeout,
		WriteTimeout:      config.HTTPServerWriteTimeout,
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
