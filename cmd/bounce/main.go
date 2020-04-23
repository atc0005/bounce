// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	htmlTemplate "html/template"
	"net/http"
	"os"
	"os/signal"
	textTemplate "text/template"

	"github.com/atc0005/bounce/config"
	"github.com/atc0005/bounce/routes"
	goteamsnotify "github.com/atc0005/go-teams-notify"
	send2teams "github.com/atc0005/send2teams/teams"

	"github.com/apex/log"
)

// for handler in cli discard es graylog json kinesis level logfmt memory
// multi papertrail text delta; do go get
// github.com/apex/log/handlers/${handler}; done

// see templates.go for the hard-coded HTML/CSS template used for the index
// page

func main() {

	// Toggle debug logging from library packages as needed to troubleshoot
	// implementation work
	goteamsnotify.DisableLogging()
	send2teams.DisableLogging()

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

	log.Debugf("AppConfig: %+v", appConfig)

	mux := http.NewServeMux()

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

	// Create context that can be used to cancel background jobs.
	ctx, cancel := context.WithCancel(context.Background())

	// Defer cancel() to cover edge cases where it might not otherwise be
	// called
	// TODO: Is this defer needed if we cover the cases elsewhere?
	defer cancel()

	// Use signal.Notify() to send a message on dedicated channel when when
	// interrupt is received (e.g., Ctrl+C) so that we can cleanly shutdown
	// the application.
	//
	// Q: Why are these channels buffered?
	// A: In order to make them asynchronous.
	// Per Bakul Shah (golang-nuts/QEORIGKZO24): In general, synchronize only
	// when you have to. Here the main thread wants to know when the worker
	// thread terminates but the worker thread doesn't care when the main
	// thread gets around to reading from "done". Using a 1 deep buffer
	// channel exactly captures this usage pattern. An unbuffered channel
	// would make the worker thread "rendezvous" with the main thread, which
	// is unnecessary.
	//
	// NOTE: Setting up a separate done channel for notify mgr and another
	// for when the http server has been shutdown.
	// done := make(chan struct{}, 1)
	httpDone := make(chan struct{}, 1)
	notifyDone := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	// Where clientRequestDetails values will be sent for processing. We use a
	// buffered channel in an effort to reduce the delay for client requests
	// as much as possible.
	notifyWorkQueue := make(chan clientRequestDetails, config.NotifyMgrQueueDepth)

	// Create "notifications manager" function that will start infinite loop
	// with select statement to process incoming notification requests.
	go StartNotifyMgr(ctx, appConfig, notifyWorkQueue, notifyDone)

	// Setup "listener" to cancel the parent contextwhen Signal.Notify()
	// indicates that SIGINT has been received
	go shutdownListener(ctx, quit, cancel)

	// Setup "listener" to shutdown the running http server when
	// the parent context has been cancelled
	go gracefulShutdown(ctx, httpServer, config.HTTPServerShutdownTimeout, httpDone)

	// Pre-process bundled templates in string/text format to Templates that
	// our handlers can execute. Based on brief testing, this seems to provide
	// a significant performance boost at the cost of a little more startup
	// time.
	indexPageHandleTemplate := htmlTemplate.Must(
		htmlTemplate.New("indexPage").Parse(handleIndexTemplateText))
	echoHandlerTemplate := textTemplate.Must(
		textTemplate.New("echoHandler").Parse(handleEchoTemplateText))

	// SETUP ROUTES
	// See handlers.go for handler definitions

	var ourRoutes routes.Routes
	ourRoutes.Add(routes.Route{
		Name:           "index",
		Description:    "Main page, fallback for unspecified routes",
		Pattern:        "/",
		AllowedMethods: []string{http.MethodGet},
		// TODO: Do we need to pass in a context here?
		HandlerFunc: handleIndex(indexPageHandleTemplate, &ourRoutes),
	})

	ourRoutes.Add(routes.Route{
		Name:           "echo",
		Description:    "Prints received values as-is to stdout and returns them via HTTP response",
		Pattern:        apiV1EchoEndpointPattern,
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
		HandlerFunc: echoHandler(
			ctx,
			echoHandlerTemplate,
			appConfig.ColorizedJSON,
			appConfig.ColorizedJSONIndent,
			notifyWorkQueue,
		),
	})

	ourRoutes.Add(routes.Route{
		Name:           "echo-json",
		Description:    "Prints formatted JSON response to stdout and via HTTP response",
		Pattern:        apiV1EchoJSONEndpointPattern,
		AllowedMethods: []string{http.MethodPost},
		HandlerFunc: echoHandler(
			ctx,
			echoHandlerTemplate,
			appConfig.ColorizedJSON,
			appConfig.ColorizedJSONIndent,
			notifyWorkQueue,
		),
	})

	ourRoutes.RegisterWithServeMux(mux)

	// listen on specified port and IP Address, block until app is terminated
	log.Infof("%s is listening on %s port %d",
		config.MyAppName, appConfig.LocalIPAddress, appConfig.LocalTCPPort)

	log.Infof("Visit http://%s:%d in your web browser for details",
		appConfig.LocalIPAddress, appConfig.LocalTCPPort)

	// TODO: This can be handled in a cleaner fashion?
	if err := httpServer.ListenAndServe(); err != nil {

		// Calling Shutdown() will immediately return ErrServerClosed, but
		// based on reading the docs it sounds like any errors from closing
		// connections will instead overwrite this default error message with
		// a real one, so receiving ErrServerClosed can be treated as a
		// "successful shutdown" message of sorts, so ignore it and look for
		// any other error message.
		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running httpServer: %v", err)
			os.Exit(1)
		}
	}

	log.Debug("Waiting on gracefulShutdown completion signal")
	<-httpDone
	log.Debug("Received gracefulShutdown completion signal")

	log.Debug("Waiting on StartNotifyMgr completion signal")
	<-notifyDone
	log.Debug("Received StartNotifyMgr completion signal")

	log.Infof("%s successfully shutdown", config.MyAppName)

}
