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
	"log"
	"net/http"
	"os"

	"github.com/atc0005/bounce/config"
	"github.com/atc0005/bounce/routes"
)

// see templates.go for the hard-coded HTML/CSS template used for the index
// page

func main() {

	log.Println("DEBUG: Initializing application")

	appConfig, err := config.NewConfig()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("Failed to initialize application: %s", err)
	}

	log.Printf("DEBUG: %+v\n", appConfig)

	// SETUP ROUTES
	// See handlers.go for handler definitions

	// TODO: replace use of http.DefaultServeMux with a custom mux to meet
	// recommended best practices

	// TODO: Use (work-in-progress) routes package to register these routes
	// for later use *and* display on the index page

	// NOTE: The entry below needs further work.
	// I need to replace renderDefaultIndexPage() (or at least update it)
	// as it currently requires arguments that don't make sense yet

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
		Description:    "Prints received values to stdout and returns them via HTTP response",
		Pattern:        apiV1EchoEndpointPattern,
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
		HandlerFunc:    echoHandler(echoHandlerTemplate),
	})

	ourRoutes.Add(routes.Route{
		Name:           "echo-json",
		Description:    "Prints formatted JSON response to stdout and via HTTP response",
		Pattern:        apiV1EchoJSONEndpointPattern,
		AllowedMethods: []string{http.MethodPost},
		HandlerFunc:    echoHandler(echoHandlerTemplate),
	})

	mux := http.NewServeMux()
	ourRoutes.RegisterWithServeMux(mux)

	// listen on specified port on ALL IP Addresses, block until app is terminated
	log.Printf("Listening on %s port %d ",
		appConfig.LocalIPAddress, appConfig.LocalTCPPort)
	listenAddress := fmt.Sprintf("%s:%d",
		appConfig.LocalIPAddress, appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}
