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

func renderDefaultIndexPage() string {

	// FIXME: Direct constant access
	return fmt.Sprintf(
		"%s\n%s\n%s",
		htmlHeader,
		htmlFallbackIndexPage,
		htmlFooter,
	)

}

func frontPageHandler() http.HandlerFunc {

	// return "type" of http.HandlerFunc as expected by http.HandleFunc() this
	// function receives `w` and `r` from http.HandleFunc; we do not have to
	// write frontPageHandler() so that it directly receives those `w` and `r`
	// as arguments.
	return func(w http.ResponseWriter, r *http.Request) {

		msgReply := fmt.Sprintf("DEBUG: frontPageHandler endpoint hit for path: %q\n", r.URL.Path)
		log.Printf(msgReply)
		//fmt.Fprintf(w, msgReply)

		// TODO: Stub out handling of non "/" requests (e.g., /favicon.ico)
		//
		// https://github.com/golang/go/issues/4799
		// https://github.com/golang/go/commit/1a819be59053fa1d6b76cb9549c9a117758090ee
		//
		// if req.URL.Path != "/" {
		// 	http.NotFound(w, req)
		// 	return
		// }

		// TODO
		// Build some kind of "banned" list?
		// Probably better to whitelist instead.
		if r.URL.Path == "/favicon.ico" {
			log.Printf("DEBUG: rejecting request for %q\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		fmt.Fprintf(w, renderDefaultIndexPage())

	}
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

	// Direct request for root of site OR unspecified route (e.g.,"catch-all")
	http.HandleFunc("/", frontPageHandler())

	// TODO: Add useful endpoints for testing here

	// listen on specified port on ALL IP Addresses, block until app is terminated
	listenAddress := fmt.Sprintf(":%d", appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
