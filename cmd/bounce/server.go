// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"
)

// shutdownListener listens for an os.Signal on the provided quit channel.
// When this signal is received, the provided parent context cancel() function
// is used to cancel all child contexts. This is intended to be run as a
// goroutine.
func shutdownListener(ctx context.Context, quit <-chan os.Signal, parentContextCancel context.CancelFunc) {

	// monitor for shutdown signal
	osSignal := <-quit

	log.Debugf("shutdownListener: Received shutdown signal: %v", osSignal)

	// Attempt to trigger a cancellation of the parent context
	log.Debug("shutdownListener: Cancelling context ...")
	parentContextCancel()
	log.Debug("shutdownListener: context canceled")

}

// gracefullShutdown listens for a context cancellation and then shuts down
// the running http server. Once the http server is shutdown, this function
// signals back that work is complete by closing the provided done channel.
// This function is intended to be run as a goroutine.
func gracefulShutdown(ctx context.Context, server *http.Server, timeout time.Duration, done chan<- struct{}) {

	log.Debug("gracefulShutdown: started; now waiting on <-ctx.Done()")

	// monitor for cancellation context
	<-ctx.Done()

	log.Debugf("gracefulShutdown: context is done: %v", ctx.Err())
	log.Warnf("%s is shutting down, please wait ...", config.MyAppName)

	// Disable HTTP keep-alives to prevent connections from persisting
	server.SetKeepAlivesEnabled(false)

	ctxShutdown, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// From the docs:
	// Shutdown returns the context's error, otherwise it returns any error
	// returned from closing the Server's underlying Listener(s).
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Errorf("gracefulShutdown: could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

// TODO: Reevaluate later once I've had a chance to think over precompilation
// of templates and whether including that step here makes any logical sense
// func newHTTPServer(serveMux *http.ServeMux) *http.Server {

// 	mux := http.NewServeMux()

// 	// Apply "default" timeout settings provided by Simon Frey; override the
// 	// default "wait forever" configuration.
// 	// FIXME: Refine these settings to apply values more appropriate for a
// 	// small-to-medium on-premise API (e.g., not over a public Internet link
// 	// where clients are expected to be slow)
// 	httpServer := http.Server{
// 		ReadHeaderTimeout: config.HTTPServerReadHeaderTimeout,
// 		ReadTimeout:       config.HTTPServerReadTimeout,
// 		WriteTimeout:      config.HTTPServerWriteTimeout,
// 		Handler:           serveMux,
// 		Addr:              fmt.Sprintf("%s:%d", appConfig.LocalIPAddress, appConfig.LocalTCPPort),
// 	}

// 	return &httpServer

// }
