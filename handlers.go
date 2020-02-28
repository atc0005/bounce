package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func frontPageHandler(w http.ResponseWriter, r *http.Request) {

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

	fmt.Fprintf(w, renderDefaultIndexPage())

}

func echoHandler(w http.ResponseWriter, r *http.Request) {

	mw := io.MultiWriter(w, os.Stdout)

	//fmt.Fprintf(w, "echoHandler endpoint hit")
	fmt.Fprintf(mw, "echoHandler endpoint hit\n")
	fmt.Fprintf(mw, "HTTP Method used by client: %s\n", r.Method)

	// https://gobyexample.com/http-servers
	fmt.Fprintf(mw, "Request received.\n")
	fmt.Fprintf(mw, "Headers:\n")

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(mw, "%v: %v\n", name, h)
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
