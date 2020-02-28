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

const (
	readme    string = "README.md"
	changelog string = "CHANGELOG.md"
)

const htmlHeader string = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

  <title>bounce - Small utility to assist with building HTTP endpoints</title>
  <meta name="description" content="bounce - Small utility to assist with building HTTP endpoints">
  <meta name="author" content="atc0005">

</head>
<body>
`

const htmlFallbackIndexPage string = `
<p>
  Welcome to the landing page for the bounce web application. This application
  is primarily intended to be used as a HTTP endpoint for testing webhook
  payloads. Over time, it may grow other related features to aid in testing
  other tools that submit data via HTTP requests.
</p>

The list of links below are the currently supported endpoints for this application:
`

const htmlFooter string = `
</body>
</html>
`

func renderDefaultIndexPage(header string, mainContent string, routes routes.Routes, footer string) string {

	// FIXME: Direct constant access
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		htmlHeader,
		htmlFallbackIndexPage,
		routes.GenerateEndPointsTable(),
		htmlFooter,
	)

}

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

	log.Printf("Listening on port %d", appConfig.LocalTCPPort)

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
		Name:        "index",
		Pattern:     "/",
		Method:      http.MethodGet,
		HandlerFunc: frontPageHandler(htmlHeader, htmlFallbackIndexPage, htmlFooter),
	})

	// Direct request for root of site OR unspecified route (e.g.,"catch-all")
	// Purpose: Landing page for list of routes, catch-all
	// http.HandleFunc("/", frontPageHandler(htmlHeader, htmlFallbackIndexPage, ourRoutes, htmlFooter))

	// // GET requests; testing endpoint
	// http.HandleFunc("/api/v1/echo", echoHandler)

	// TODO: Add useful endpoints for testing here

	// listen on specified port on ALL IP Addresses, block until app is terminated
	listenAddress := fmt.Sprintf(":%d", appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
