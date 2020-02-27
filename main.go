// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/atc0005/bounce/config"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func loadMarkdown(filename string) ([]byte, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("error loading Markdown file %q: %s", filename, err)
		return nil, err
	}
	return data, nil
}

// processMarkdown runs a Markdown processor against the stored Page content
// and replaces supported Markdown with HTML equivalents for display to
// the client.
func processMarkdown(b []byte, skipSanitize bool) ([]byte, error) {

	// add protection against nil pointer deference
	if b == nil {
		return nil, fmt.Errorf("aborting processing of nil pointer")
	}

	if !skipSanitize {
		unsafe := blackfriday.Run(b)
		data := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		return data, nil
	}

	data := blackfriday.Run(b)

	return data, nil

}

func frontPageHandler(skipSanitize bool) http.HandlerFunc {

	// return "type" of http.HandlerFunc as expected by http.HandleFunc() this
	// function receives `w` and `r` from http.HandleFunc; we do not have to
	// write frontPageHandler() so that it directly receives those `w` and `r`
	// as arguments.
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("frontPageHandler endpoint hit")
		//fmt.Fprintf(w, "frontPageHandler endpoint hit")

		filename := "README.md"
		markdownInput, err := loadMarkdown(filename)
		if err != nil {
			log.Fatalf("Unable to open %s: %s", filename, err)
		}
		bytes, err := processMarkdown(markdownInput, skipSanitize)
		htmlOutput := string(bytes)
		fmt.Fprintf(w, htmlOutput)

	}
}

func main() {

	log.Println("Initializing application")

	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize application: %s", err)
	}

	log.Printf("%+v\n", appConfig)

	log.Printf("Listening on port %d", appConfig.LocalTCPPort)

	http.HandleFunc("/", frontPageHandler(appConfig.SkipMarkdownSanitization))

	// listen on port 8080 on any interface, block until app is terminated
	log.Fatal(http.ListenAndServe(":8000", nil))
}
