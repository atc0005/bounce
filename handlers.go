// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/atc0005/bounce/routes"
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

func echoHandler(w http.ResponseWriter, r *http.Request) {

	mw := io.MultiWriter(w, os.Stdout)

	//fmt.Fprintf(w, "echoHandler endpoint hit")
	fmt.Fprintf(mw, "DEBUG: echoHandler endpoint hit\n\n")

	fmt.Fprintf(mw, "HTTP Method used by client: %s\n", r.Method)
	fmt.Fprintf(mw, "Client IP Address: %s\n", GetIP(r))

	fmt.Fprintf(mw, "\nHeaders:\n\n")

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(mw, "  * %v: %v\n", name, h)
		}
	}

	// Only try to get the body if the client submitted a payload
	if r.Method == http.MethodPost {
		fmt.Fprintf(mw, "POST request received; reading Body value ...\n")

		fmt.Fprintf(mw, "Body:\n")
		_, err := io.Copy(mw, r.Body)
		if err != nil {
			log.Println(err)
			return
		}
	}

}
