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

	// Direct request for root of site OR unspecified route (e.g.,"catch-all")
	// Purpose: Landing page for list of routes, catch-all
	http.HandleFunc("/", frontPageHandler)

	// GET requests; testing endpoint
	http.HandleFunc("/api/v1/echo", echoHandler)

	// TODO: Add useful endpoints for testing here

	// https://stackoverflow.com/a/37244145/903870
	// v := reflect.ValueOf(http.DefaultServeMux).Elem()
	// fmt.Printf("routes: %v\n", v.FieldByName("m"))

	// https://stackoverflow.com/a/37247014/903870
	// httpMux := reflect.ValueOf(http.DefaultServeMux).Elem()
	// finList := httpMux.FieldByIndex([]int{1})
	// fmt.Println(finList)

	// Broken attempt below
	//
	// httpMux := reflect.ValueOf(http.DefaultServeMux).Elem()
	// finList := httpMux.FieldByIndex([]int{1})
	// fmt.Println(finList)
	// fmt.Printf("finList is of type: %T", finList)
	// var routes map[string]string
	// var ok bool
	// if routes, ok = finList.Interface().(map[string]string); !ok {
	// 	fmt.Println("Unable to properly obtain list of routes")
	// }
	// fmt.Println("Defined routes:")
	// for k, v := range routes {
	// 	fmt.Println("Key:", k, "Value:", v)
	// }

	// 	type MyInt int
	// var x MyInt = 7
	// v := reflect.ValueOf(x)
	// y := v.Interface().(float64) // y will have type float64.
	// fmt.Println(y)

	// listen on specified port on ALL IP Addresses, block until app is terminated
	listenAddress := fmt.Sprintf(":%d", appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
