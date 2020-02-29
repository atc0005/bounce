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

const htmlTemplate string = `
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">

  <title>bounce - Small utility to assist with building HTTP endpoints</title>
  <meta name="description" content="bounce - Small utility to assist with building HTTP endpoints">
  <meta name="author" content="atc0005">

  <!--
	  https://www.w3schools.com/css/css_table.asp
	  http://web.simmons.edu/~grabiner/comm244/weekfour/code-test.html
  -->
  <style>
  table {
	border-collapse: collapse;
	width: 100%;
  }

  th, td {
	text-align: left;
	padding: 8px;
  }

  tr:nth-child(even){background-color: #f2f2f2}

  th {
	background-color: #4CAF50;
	color: white;
  }

  code {
	background-color: #eee;
	border: 1px solid #999;
	display: block;
	padding: 0.5em;
  }

  </style>

</head>
<body>

<h1>Welcome!</h1>

<p>
  Welcome to the landing page for the bounce web application. This application
  is primarily intended to be used as a HTTP endpoint for testing webhook
  payloads. Over time, it may grow other related features to aid in testing
  other tools that submit data via HTTP requests.
</p>

<p>
  The list of links below are the currently supported endpoints for this
  application:
</p>

<table>
  <tr>
    <th>Name</th>
    <th>Pattern</th>
    <th>Description</th>
    <th>Allowed Methods</th>
  </tr>
{{range .}}
  <tr>
    <td><code>{{ .Name }}</code></td>
    <td><code>{{ .Pattern }}</code></td>
	<td><code>{{ .Description }}</td>
	<td>{{range .AllowedMethods}}<code>{{ . }}</code> {{end}}</td>
  </tr>
{{else}}
<tr>
  <td><code>Failed to parse routes</code></td>
  <td><code>N/A</code></td>
  <td><code>N/A</code></td>
  <td><code>N/A</code></td>
</tr>
{{end}}
</table>

</body>
</html>
`

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
		HandlerFunc:    handleIndex(htmlTemplate, &ourRoutes),
	})

	ourRoutes.Add(routes.Route{
		Name:           "echo",
		Description:    "The echo endpoint prints received values to stdout and returns them via HTTP response",
		Pattern:        "/echo",
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
		HandlerFunc:    echoHandler,
	})

	mux := http.NewServeMux()
	ourRoutes.RegisterWithServeMux(mux)

	// listen on specified port on ALL IP Addresses, block until app is terminated
	log.Printf("Listening on port %d", appConfig.LocalTCPPort)
	listenAddress := fmt.Sprintf(":%d", appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}
