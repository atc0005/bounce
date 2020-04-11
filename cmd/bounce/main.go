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
	textTemplate "text/template"
	"time"

	"github.com/atc0005/bounce/config"
	"github.com/atc0005/bounce/routes"
	goteamsnotify "github.com/atc0005/go-teams-notify"
	send2teams "github.com/atc0005/send2teams/teams"

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

func init() {

	// Go ahead and enable debug logging from these library packages while we
	// are actively working on the `i21-add-msteams-integration-2nd-attempt`
	// branch
	// goteamsnotify.EnableLogging()
	// send2teams.EnableLogging()

	goteamsnotify.EnableLogging()
	send2teams.DisableLogging()
}

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

	// Create context that can be used to cancel background jobs.
	ctx, cancel := context.WithCancel(context.Background())

	// cancel when we are finished sending notification requests
	defer cancel()

	// Where echoHandlerResponse values will be sent for processing. We use a
	// buffered channel in an effort to reduce the delay for client requests
	// as much as possible.
	notifyWorkQueue := make(chan echoHandlerResponse, 5)

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

	// Create "notifications manager" function that will start infinite loop
	// with select statement to process incoming notification requests.
	go StartNotifyMgr(ctx, appConfig, notifyWorkQueue)

	// listen on specified port on ALL IP Addresses, block until app is terminated
	log.Infof("Listening on %s port %d ",
		appConfig.LocalIPAddress, appConfig.LocalTCPPort)

	// FIXME: Remove this once done testing
	// TODO: Replace with use of Signal to call cancel() when Ctrl+C is issued?
	go func() {
		time.Sleep(time.Second * 3)
		log.Warn("Calling cancel() to test shutdown behavior for notifier")
		cancel()
	}()

	// Setup "listener" to shutdown the running http server when a
	// cancellation context is triggered
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():

			ctxErr := ctx.Err()

			log.Debugf("main: Received Done signal: %v, shutting down ...", ctxErr)

			ctxShutDown, cancel := context.WithTimeout(ctx, 5*time.Second)

			// what is this cancelling exactly?
			defer cancel()

			// Pass in a new timeout-based context to *force* shutdown if the
			// normal shutdown process takes longer than expected.
			err := httpServer.Shutdown(ctxShutDown)
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Errorf("main: error shutting down http server: %v", err)
			}
		}
	}(ctx)

	// TODO: This can be handled in a cleaner fashion?
	if err := httpServer.ListenAndServe(); err != nil {

		if !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running httpServer: %v", err)
		} else {

			// Calling Shutdown() will immediately return ErrServerClosed, but
			// based on reading the docs it sounds like any errors from
			// closing connections will instead overwrite this default error
			// message with a real one, so receiving ErrServerClosed can be
			// treated as a "successful shutdown" message of sorts.
			log.Debug("main: successfully shutdown httpServer")
		}

		// the deferred cancel() from earlier should be sufficient to handle
		// this task, but we call it explicitly just to be sure.
		// TODO: Is this best practice? Is it safe to call cancel() multiple
		// times?
		log.Debug("Explicitly using cancel() to shutdown background tasks")
		cancel()
	}

}
