// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/atc0005/bounce/routes"

	"github.com/golang/gddo/httputil/header"
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

		t := time.Now()

		//fmt.Fprintf(w, "echoHandler endpoint hit")
		fmt.Fprintf(mw, "DEBUG: echoHandler endpoint hit\n\n")
		fmt.Fprintf(mw, "Request received: %v\n", t.Format(time.RFC3339))

		fmt.Fprintf(mw, "Endpoint path requested by client: %s\n", r.URL.Path)
		fmt.Fprintf(mw, "HTTP Method used by client: %s\n", r.Method)
		fmt.Fprintf(mw, "Client IP Address: %s\n", GetIP(r))

		fmt.Fprintf(mw, "\nHeaders:\n\n")

		for name, headers := range r.Header {
			for _, h := range headers {
				fmt.Fprintf(mw, "  * %v: %v\n", name, h)
			}
		}

		fmt.Fprintf(mw, "\nUnformatted Body:\n\n")

		switch r.Method {

		case http.MethodGet:
			fmt.Fprintf(mw, "Sorry, this endpoint only accepts JSON data via %s requests.\n", http.MethodPost)
			fmt.Fprintf(mw, "Please see the README for examples and then try again.\n")

		case http.MethodPost:

			var err error

			// /api/v1/echo
			// /api/v1/echo/json
			//
			//
			// return 404 if not one of those EXACT endpoints
			// return helpful text if /api/v1/echo/json and NOT POST method or expected content-type

			// Copy body to a buffer since we'll use it in multiple places and
			// (I think?) you can only read from r.Body once
			// buffer := bytes.Buffer{}
			// _, err = io.Copy(&buffer, r.Body)
			// requestBodyReader := bytes.NewReader(buffer.Bytes())
			requestBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errorMsg := fmt.Sprintf("Error reading request body: %s", err)
				fmt.Fprintf(mw, errorMsg)
				http.Error(w, errorMsg, http.StatusBadRequest)
				return
			}
			requestBodyBuffer := bytes.NewBuffer(requestBody)
			// TODO: Do we really need an io.ReadCloser here?
			requestBodyReader := ioutil.NopCloser(requestBodyBuffer)

			// write out request body in raw format
			_, err = io.Copy(mw, requestBodyReader)
			if err != nil {
				log.Println(err)
				return
			}

			// Add some whitespace to separate previous/upcoming contents
			fmt.Fprintf(mw, "\n\n")

			// Only attempt to parse the request body as JSON if the
			// JSON-specific endpoint was used
			if r.URL.Path == apiV1EchoJSONEndpointPattern {

				// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
				//
				// If the Content-Type header is present, check that it has the value
				// application/json. Note that we are using the gddo/httputil/header
				// package to parse and extract the value here, so the check works
				// even if the client includes additional charset or boundary
				// information in the header.
				contentTypeHeader := r.Header.Get("Content-Type")
				if contentTypeHeader != "" {
					value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
					if value != "application/json" {
						msg := fmt.Sprintf("Submitted request %q does not contain the expected application/json Content-Type header.", contentTypeHeader)
						fmt.Fprintf(os.Stdout, msg)
						http.Error(w, msg, http.StatusUnsupportedMediaType)
						return
					}
				}

				fmt.Fprintf(mw, "Formatted Body:\n")

				// https://golang.org/pkg/encoding/json/#Indent
				var prettyJSON bytes.Buffer
				// FIXME: Is it safe now to access requestBody (byte slice)
				// directly with all of the additional "wrappers" applied to it?
				err = json.Indent(&prettyJSON, requestBody, "", "\t")
				if err != nil {
					errorMsg := fmt.Sprintf("JSON parse error: %s", err)
					fmt.Fprintf(os.Stdout, errorMsg)
					http.Error(w, errorMsg, http.StatusBadRequest)
					return
				}
				fmt.Fprintf(mw, prettyJSON.String())

				// https://golang.org/pkg/encoding/json/#MarshalIndent
				// prettyJSON, err := json.MarshalIndent(&buffer, "", "\t")
				// if err != nil {
				// 	errorMsg := fmt.Sprintf("JSON parse error: %s", err)
				// 	fmt.Fprintf(mw, errorMsg)
				// 	http.Error(w, errorMsg, http.StatusBadRequest)
				// 	return
				// }
				// fmt.Fprintf(mw, string(prettyJSON))

			}

			// Add some whitespace to separate previous/upcoming contents
			fmt.Fprintf(mw, "\n\n\n")

		default:
			fmt.Fprintf(mw, "ERROR: Unsupported method %q received; please try again using %s method\n", r.Method, http.MethodPost)

		}

	default:
		log.Printf("DEBUG: Rejecting request %q; not explicitly handled by a route.\n", r.URL.Path)
		http.NotFound(w, r)
		return
	}

}
