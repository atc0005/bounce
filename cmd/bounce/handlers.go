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
	"errors"
	"fmt"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	textTemplate "text/template"
	"time"

	"github.com/atc0005/bounce/routes"

	"github.com/TylerBrock/colorjson"
	"github.com/apex/log"
	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
)

// API endpoint patterns supported by this application
//
// TODO: Find a better location for these values
const (
	apiV1EchoEndpointPattern     string = "/api/v1/echo"
	apiV1EchoJSONEndpointPattern string = "/api/v1/echo/json"
)

// MB is a convenience constant that represents 1 Megabyte, which so happens
// to be the limit we're placing on request bodies (in order to help limit the
// impact from misconfigured http clients).
// TODO: Find a better location for this constant
const MB int64 = 1048576

// echoHandlerResponse is used to bundle various client request details for
// processing by templates or notification functions.
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

// handleIndex receives our HTML template and our defined routes as a pointer.
// Both are used to generate a dynamic index of the available routes or
// "endpoints" for users to target with test payloads. A pointer is used because
// by the time this handler is defined, the full set of routes has *not* been
// defined. Using a pointer, we are able to access the complete collection
// of defined routes when this handler is finally called.
func handleIndex(templateText string, rs *routes.Routes) http.HandlerFunc {

	tmpl := htmlTemplate.Must(htmlTemplate.New("indexPage").Parse(templateText))

	return func(w http.ResponseWriter, r *http.Request) {

		ctxLog := log.WithFields(log.Fields{
			"url_path":   r.URL.Path,
			"num_routes": len(*rs),
		})

		ctxLog.Debug("handleIndex endpoint hit")

		if r.Method != http.MethodGet {

			ctxLog.WithFields(log.Fields{
				"http_method": r.Method,
			}).Debug("non-GET request received on GET-only endpoint")
			errorMsg := fmt.Sprintf(
				"Sorry, this endpoint only accepts %s requests. "+
					"Please see the README for examples and then try again.",
				http.MethodGet,
			)
			// TODO: Can apex/log hook into this and handle output?
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			fmt.Fprint(w, errorMsg)
			return
		}

		// https://github.com/golang/go/issues/4799
		// https://github.com/golang/go/commit/1a819be59053fa1d6b76cb9549c9a117758090ee
		if r.URL.Path != "/" {
			ctxLog.Debug("Rejecting request not explicitly handled by a route")
			http.NotFound(w, r)
			return
		}

		for _, route := range *rs {
			log.Debugf("route: %+v", route)
		}

		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, *rs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ctxLog.Error(err.Error())
		}

	}

}

// echoHandler echos back the HTTP request received by
func echoHandler(templateText string, coloredJSON bool, coloredJSONIndent int, webhookURL string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// For now, we generate plain text responses
		//w.Header().Set("Content-Type", "text/plain")

		ourResponse := echoHandlerResponse{}

		mw := io.MultiWriter(w, os.Stdout)

		tmpl := textTemplate.Must(textTemplate.New("echoHandler").Parse(templateText))

		writeTemplate := func() {
			err := tmpl.Execute(mw, ourResponse)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Errorf("error occurred while trying to execute template: %v", err)

				// We force a a return here since it is unlikely that we should
				// execute any other code after failing to generate/write out our
				// template
				return
			}

			// Manually flush http.ResponseWriter
			// https://blog.simon-frey.eu/manual-flush-golang-http-responsewriter/
			if f, ok := w.(http.Flusher); ok {
				log.Debug("Manually flushing http.ResponseWriter")
				f.Flush()
			} else {
				log.Warn("http.Flusher interface not available, cannot flush http.ResponseWriter")
				log.Warn("Not flushing http.ResponseWriter may cause a noticeable delay between requests")
			}

		}

		log.Debug("echoHandler endpoint hit")

		ourResponse.Datestamp = time.Now().Format((time.RFC3339))
		ourResponse.EndpointPath = r.URL.Path
		ourResponse.HTTPMethod = r.Method
		ourResponse.ClientIPAddress = GetIP(r)
		ourResponse.Headers = r.Header

		switch r.URL.Path {

		// Expected endpoint patterns for this handler
		case apiV1EchoEndpointPattern:

			switch r.Method {

			case http.MethodGet:

				// Write out what we have.
				writeTemplate()

				// Send request details to Microsoft Teams if webhook URL set
				if webhookURL != "" {
					ourMessage := createMessage(ourResponse)
					if err := sendMessage(webhookURL, ourMessage); err != nil {
						log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
					}
				}

				return

			case http.MethodPost:

				// Limit request body to 1 MB
				r.Body = http.MaxBytesReader(w, r.Body, 1*MB)
				requestBody, err := ioutil.ReadAll(r.Body)
				if err != nil {
					errorMsg := fmt.Sprintf("Error reading request body: %s", err)
					ourResponse.BodyError = errorMsg

					http.Error(w, errorMsg, http.StatusBadRequest)
					log.Error(errorMsg)

					writeTemplate()

					// Send request details to Microsoft Teams if webhook URL set
					if webhookURL != "" {
						ourMessage := createMessage(ourResponse)
						if err := sendMessage(webhookURL, ourMessage); err != nil {
							log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
						}
					}

					return
				}

				ourResponse.Body = string(requestBody)
				ourResponse.FormattedBodyError = fmt.Sprintf(
					"This endpoint does not apply JSON formatting to the request body.\n"+
						"Use the %q endpoint for JSON payload testing.",
					apiV1EchoJSONEndpointPattern,
				)

				// If we made it this far, then presumably our template data
				// structure "ourResponse" is fully populated and we can execute
				// the template against it
				writeTemplate()

				// Send request details to Microsoft Teams if webhook URL set
				if webhookURL != "" {
					ourMessage := createMessage(ourResponse)
					if err := sendMessage(webhookURL, ourMessage); err != nil {
						log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
					}
				}

				return

			default:
				errorMsg := fmt.Sprintf("ERROR: Unsupported method %q received; please try again using %s method", r.Method, http.MethodPost)
				ourResponse.RequestError = errorMsg

				http.Error(w, errorMsg, http.StatusMethodNotAllowed)
				log.Error(errorMsg)

				writeTemplate()

				// Send request details to Microsoft Teams if webhook URL set
				if webhookURL != "" {
					ourMessage := createMessage(ourResponse)
					if err := sendMessage(webhookURL, ourMessage); err != nil {
						log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
					}
				}

				return
			}

		case apiV1EchoJSONEndpointPattern:

			switch r.Method {

			case http.MethodGet:
				// TODO: Collect this for use with our template
				errorMsg := fmt.Sprintf(
					"Sorry, this endpoint only accepts JSON data via %s requests. "+
						"Please see the README for examples and then try again.",
					http.MethodPost,
				)
				ourResponse.RequestError = errorMsg

				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				log.Error(errorMsg)

				writeTemplate()

				// Send request details to Microsoft Teams if webhook URL set
				if webhookURL != "" {
					ourMessage := createMessage(ourResponse)
					if err := sendMessage(webhookURL, ourMessage); err != nil {
						log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
					}
				}

				return

			case http.MethodPost:

				// Limit request body to 1 MB
				r.Body = http.MaxBytesReader(w, r.Body, 1*MB)

				// read everything from the (size-limited) request body so
				// that we can display it in a raw format, replace the Body
				// with a new io.ReadCloser to allow later access to r.Body
				// for JSON-decoding purposes
				requestBody, err := ioutil.ReadAll(r.Body)
				r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

				if err != nil {
					errorMsg := fmt.Sprintf("Error reading request body: %s", err)
					ourResponse.BodyError = errorMsg

					http.Error(w, errorMsg, http.StatusBadRequest)
					log.Error(errorMsg)

					writeTemplate()

					// Send request details to Microsoft Teams if webhook URL set
					if webhookURL != "" {
						ourMessage := createMessage(ourResponse)
						if err := sendMessage(webhookURL, ourMessage); err != nil {
							log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
						}
					}

					return
				}
				ourResponse.Body = string(requestBody)

				handleJSONParseError := func(w http.ResponseWriter, err error) {
					if err != nil {

						var mr *malformedRequest
						errorPrefix := "JSON parse error"
						if errors.As(err, &mr) {
							log.WithFields(log.Fields{
								"msg":    mr.msg,
								"status": mr.status,
							}).Error(errorPrefix)

							ourResponse.FormattedBodyError = fmt.Sprintf("%s: %s", errorPrefix, mr.msg)

							http.Error(w, mr.msg, mr.status)

							writeTemplate()

							// Send request details to Microsoft Teams if webhook URL set
							if webhookURL != "" {
								ourMessage := createMessage(ourResponse)
								if err := sendMessage(webhookURL, ourMessage); err != nil {
									log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
								}
							}

							return
						}

						errorMsg := fmt.Sprintf("%s: %s", errorPrefix, err.Error())
						ourResponse.FormattedBodyError = errorMsg
						http.Error(w, errorMsg, http.StatusInternalServerError)
						log.Error(errorMsg)

						writeTemplate()

						// Send request details to Microsoft Teams if webhook URL set
						if webhookURL != "" {
							ourMessage := createMessage(ourResponse)
							if err := sendMessage(webhookURL, ourMessage); err != nil {
								log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
							}
						}

						return
					}
				}

				// Decode request body into JSON using helper function
				var decodedJSON map[string]interface{}

				// At this point we're dealing with a `malformedRequest` type
				// of error. We can reference recorded `status` and `msg`
				// fields to provide more information. Our
				// `handleJSONParseError()` helper function looks for this
				// type and uses it as that type if found.
				err = decodeJSONBody(w, r, &decodedJSON)
				handleJSONParseError(w, err)

				switch coloredJSON {
				case true:
					// Make a custom formatter with indent set
					colorJSONFormatter := colorjson.NewFormatter()
					colorJSONFormatter.Indent = coloredJSONIndent

					// Marshall into Colorized JSON
					jsonBytes, err := colorJSONFormatter.Marshal(decodedJSON)
					handleJSONParseError(w, err)
					ourResponse.FormattedBody = string(jsonBytes)

				case false:
					// https://golang.org/pkg/encoding/json/#Indent
					var prettyJSON bytes.Buffer
					err = json.Indent(&prettyJSON, requestBody, "", "\t")
					handleJSONParseError(w, err)
					ourResponse.FormattedBody = prettyJSON.String()
				}

				// If we made it this far, then presumably our template data
				// structure "ourResponse" is fully populated and we can execute
				// the template against it
				writeTemplate()

				// Send request details to Microsoft Teams if webhook URL set
				if webhookURL != "" {
					ourMessage := createMessage(ourResponse)
					if err := sendMessage(webhookURL, ourMessage); err != nil {
						log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
					}
				}

			default:
				errorMsg := fmt.Sprintf("ERROR: Unsupported method %q received; please try again using %s method", r.Method, http.MethodPost)
				ourResponse.RequestError = errorMsg

				http.Error(w, errorMsg, http.StatusMethodNotAllowed)

				writeTemplate()

				// Send request details to Microsoft Teams if webhook URL set
				if webhookURL != "" {
					ourMessage := createMessage(ourResponse)
					if err := sendMessage(webhookURL, ourMessage); err != nil {
						log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
					}
				}

				return
			}

		default:
			// Template is not used for this code block, so no need to account for
			// the output in the template
			log.Debugf("Rejecting request %q; not explicitly handled by a route.", r.URL.Path)
			http.NotFound(w, r)

			// Send request details to Microsoft Teams if webhook URL set
			if webhookURL != "" {
				ourMessage := createMessage(ourResponse)
				if err := sendMessage(webhookURL, ourMessage); err != nil {
					log.Errorf("error occurred while trying to send message to Microsoft Teams: %v", err)
				}
			}

			return
		}

	}
}
