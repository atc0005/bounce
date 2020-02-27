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
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/atc0005/bounce/config"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
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

const htmlFooter string = `
</body>
</html>
`

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

func frontPageHandler(requestedFile string, skipSanitize bool) http.HandlerFunc {

	// return "type" of http.HandlerFunc as expected by http.HandleFunc() this
	// function receives `w` and `r` from http.HandleFunc; we do not have to
	// write frontPageHandler() so that it directly receives those `w` and `r`
	// as arguments.
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("DEBUG: frontPageHandler endpoint hit")
		//fmt.Fprintf(w, "frontPageHandler endpoint hit")

		filename := requestedFile
		markdownInput, err := loadMarkdown(filename)
		if err != nil {
			log.Fatalf("Unable to open %s: %s", filename, err)
		}
		bytes, err := processMarkdown(markdownInput, skipSanitize)
		htmlOutput := fmt.Sprintf(
			"%s\n%s\n%s",
			htmlHeader,
			htmlFooter,
			string(bytes),
		)
		fmt.Fprintf(w, htmlOutput)

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

	// Setup routes
	http.HandleFunc("/", frontPageHandler(readme, appConfig.SkipMarkdownSanitization))
	http.HandleFunc(fmt.Sprintf("/%s", changelog), frontPageHandler(changelog, appConfig.SkipMarkdownSanitization))

	// listen on specified port on any interface, block until app is terminated
	listenAddress := fmt.Sprintf(":%d", appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
