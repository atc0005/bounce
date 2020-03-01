// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/atc0005/bounce/routes"
)

// API endpoint patterns supported by this application
//
// TODO: Find a better location for these values
const (
	apiV1EchoEndpointPattern     string = "/api/v1/echo"
	apiV1EchoJSONEndpointPattern string = "/api/v1/echo/json"
)

// handleIndex receives our HTML template and our defined routes as a pointer.
// Both are used to generate a dynamic index of the available routes or
// "endpoints" for users to target with test payloads. A pointer is used because
// by the time this handler is defined, the full set of routes has *not* been
// defined. Using a pointer, we are able to access the complete collection
// of defined routes when this handler is finally called.
func handleIndex(htmlTemplateText string, rs *routes.Routes) http.HandlerFunc {

	htmlTemplate := template.Must(template.New("indexPage").Parse(htmlTemplate))

	// FIXME: Guard against POST requests to this endpoint?

	return func(w http.ResponseWriter, r *http.Request) {

		msgReply := fmt.Sprintf("DEBUG: handleIndex endpoint hit for path: %q\n", r.URL.Path)
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
		// if r.URL.Path == "/favicon.ico" {
		// 	log.Printf("DEBUG: rejecting request for %q\n", r.URL.Path)
		// 	http.NotFound(w, r)
		// 	return
		// }

		if r.URL.Path != "/" {
			log.Printf("DEBUG: Rejecting request %q; not explicitly handled by a route.\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		log.Println("DEBUG: length of routes:", len(*rs))
		log.Printf("DEBUG: Type of *rs is %T. Fields: %v+", *rs, *rs)

		for _, route := range *rs {
			log.Println("DEBUG: route:", route)
		}

		htmlTemplate.Execute(w, *rs)

	}

}

// TODO: Convert this so that it serves multiple endpoints
//
// /api/v1/echo
// /api/v1/echo/json
//
//
// return 404 if not one of those EXACT endpoints
// return helpful text if /api/v1/echo/json and NOT POST method or expected content-type
func echoHandler(w http.ResponseWriter, r *http.Request) {

	switch r.URL.Path {

	// Expected endpoint patterns for this handler
	case apiV1EchoEndpointPattern, apiV1EchoJSONEndpointPattern:

		mw := io.MultiWriter(w, os.Stdout)

		//fmt.Fprintf(w, "echoHandler endpoint hit")
		fmt.Fprintf(mw, "DEBUG: echoHandler endpoint hit\n\n")

		fmt.Fprintf(mw, "Endpoint path requested by client: %s\n", r.URL.Path)
		fmt.Fprintf(mw, "HTTP Method used by client: %s\n", r.Method)
		fmt.Fprintf(mw, "Client IP Address: %s\n", GetIP(r))

		fmt.Fprintf(mw, "\nHeaders:\n\n")

		for name, headers := range r.Header {
			for _, h := range headers {
				fmt.Fprintf(mw, "  * %v: %v\n", name, h)
			}
		}

		switch r.Method {

		case http.MethodGet:
			fmt.Fprintf(mw, "Sorry, this endpoint only accepts JSON data via %s requests.\n", http.MethodPost)
			fmt.Fprintf(mw, "Please see the README for examples and then try again.\n")

		case http.MethodPost:

			// /api/v1/echo
			// /api/v1/echo/json
			//
			//
			// return 404 if not one of those EXACT endpoints
			// return helpful text if /api/v1/echo/json and NOT POST method or expected content-type

			// Copy body to a buffer since we'll use it in multiple places and
			// (I think?) you can only read from r.Body once
			buffer := bytes.Buffer{}
			_, err := io.Copy(&buffer, r.Body)

			fmt.Fprintf(mw, "Body:\n")
			_, err = io.Copy(mw, &buffer)
			if err != nil {
				log.Println(err)
				return
			}

		default:
			fmt.Fprintf(mw, "ERROR: Unsupported method %q received; please try again using %s method\n", r.Method, http.MethodPost)

		}

	default:
		log.Printf("DEBUG: Rejecting request %q; not explicitly handled by a route.\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}

}
