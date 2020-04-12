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

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"
)

// gracefullShutdown runs is intended to run as a goroutine and listen for a
// shutdown signal. When this signal is received on a provided quit channel,
// the http server is shutdown. Once the http server is shutdown, this
// function signals back that work is complete by closing the provided done
// channel.
func gracefulShutdown(ctx context.Context, server *http.Server, quit <-chan os.Signal, done chan<- struct{}) {

	// monitor for shutdown signal
	osSignal := <-quit

	log.Warnf("system call:%+v", osSignal)
	log.Warn("Received shutdown signal: %v")
	log.Warn("Server is shutting down, please wait ...")

	// Disable HTTP keep-alives to prevent connections from persisting
	server.SetKeepAlivesEnabled(false)

	ctxShutdown, cancel := context.WithTimeout(ctx, config.HTTPServerShutdownTimeout)
	defer cancel()

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
