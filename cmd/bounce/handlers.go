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
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	textTemplate "text/template"
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
func handleIndex(templateText string, rs *routes.Routes) http.HandlerFunc {

	htmlTemplate := htmlTemplate.Must(htmlTemplate.New("indexPage").Parse(templateText))

	// FIXME: Guard against POST requests to this endpoint?

	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("DEBUG: handleIndex endpoint hit for path: %q\n", r.URL.Path)

		if r.Method != http.MethodGet {

			log.Println("DEBUG: non-GET request received on GET-only endpoint")
			errorMsg := fmt.Sprintf(
				"Sorry, this endpoint only accepts %s requests.\n"+
					"Please see the README for examples and then try again.",
				http.MethodGet,
			)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Fprint(w, errorMsg)
			return
		}

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

		w.Header().Set("Content-Type", "text/html")
		err := htmlTemplate.Execute(w, *rs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

}

// echoHandler echos back the HTTP request received by
func echoHandler(templateText string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// For now, we generate plain text responses
		//w.Header().Set("Content-Type", "text/plain")

		type echoHandlerResponse struct {
			Datestamp          string
			EndpointPath       string
			HTTPMethod         string
			ClientIPAddress    string
			Headers            http.Header
			Body               string
			BodyError          string
			FormattedBody      string
			FormattedBodyError string
			RequestError       string
			ContentTypeError   string
		}

		ourResponse := echoHandlerResponse{}

		mw := io.MultiWriter(w, os.Stdout)

		textTemplate := textTemplate.Must(textTemplate.New("echoHandler").Parse(templateText))

		writeTemplate := func() {
			err := textTemplate.Execute(mw, ourResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				// We force a a return here since it is unlikely that we should
				// execute any other code after failing to generate/write out our
				// template
				return
			}
		}

		fmt.Fprintf(os.Stdout, "DEBUG: echoHandler endpoint hit\n\n")

		ourResponse.Datestamp = time.Now().Format((time.RFC3339))
		ourResponse.EndpointPath = r.URL.Path
		ourResponse.HTTPMethod = r.Method
		ourResponse.ClientIPAddress = GetIP(r)
		ourResponse.Headers = r.Header

		switch r.URL.Path {

		// Expected endpoint patterns for this handler
		case apiV1EchoEndpointPattern:

			// Write out what we have. Nothing further to do for this endpoint
			writeTemplate()
			return

		case apiV1EchoJSONEndpointPattern:

			switch r.Method {

			case http.MethodGet:
				// TODO: Collect this for use with our template
				errorMsg := fmt.Sprintf(
					"Sorry, this endpoint only accepts JSON data via %s requests.\n"+
						"Please see the README for examples and then try again.",
					http.MethodPost,
				)
				ourResponse.RequestError = errorMsg

				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				writeTemplate()
				return

			case http.MethodPost:

				var err error

				// Copy body to a buffer since we'll use it in multiple places and
				// (I think?) you can only read from r.Body once
				// buffer := bytes.Buffer{}
				// _, err = io.Copy(&buffer, r.Body)
				// requestBodyReader := bytes.NewReader(buffer.Bytes())
				requestBody, err := ioutil.ReadAll(r.Body)
				if err != nil {
					errorMsg := fmt.Sprintf("Error reading request body: %s", err)
					ourResponse.BodyError = errorMsg

					http.Error(w, errorMsg, http.StatusBadRequest)
					writeTemplate()
					return
				}
				//requestBodyBuffer := bytes.NewBuffer(requestBody)
				// TODO: Do we really need an io.ReadCloser here?
				//requestBodyReader := ioutil.NopCloser(requestBodyBuffer)

				ourResponse.Body = string(requestBody)

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
							errorMsg := fmt.Sprintf("Submitted request %q does not contain the expected application/json Content-Type header.", contentTypeHeader)
							ourResponse.ContentTypeError = errorMsg

							http.Error(w, errorMsg, http.StatusUnsupportedMediaType)
							writeTemplate()
							return
						}
					}

					// https://golang.org/pkg/encoding/json/#Indent
					var prettyJSON bytes.Buffer
					// FIXME: Is it safe now to access requestBody (byte slice)
					// directly with all of the additional "wrappers" applied to it?
					err = json.Indent(&prettyJSON, requestBody, "", "\t")
					if err != nil {
						errorMsg := fmt.Sprintf("JSON parse error: %s", err)
						ourResponse.FormattedBodyError = errorMsg

						http.Error(w, errorMsg, http.StatusBadRequest)
						writeTemplate()
						return
					}
					ourResponse.FormattedBody = prettyJSON.String()

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

				// If we made it this far, then presumably our template data
				// structure "ourResponse" is fully populated and we can execute
				// the template against it
				writeTemplate()

			default:
				errorMsg := fmt.Sprintf("ERROR: Unsupported method %q received; please try again using %s method", r.Method, http.MethodPost)
				ourResponse.RequestError = errorMsg

				http.Error(w, errorMsg, http.StatusMethodNotAllowed)
				writeTemplate()
				return
			}

		default:
			// Template is not used for this code block, so no need to account for
			// the output in the template
			log.Printf("DEBUG: Rejecting request %q; not explicitly handled by a route.\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}

	}
}
