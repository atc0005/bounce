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

	"github.com/shurcooL/github_flavored_markdown"
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

func loadMarkdown(filename string) ([]byte, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("DEBUG: error loading Markdown file %q: %s", filename, err)
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

	if skipSanitize {
		log.Printf("DEBUG: Skipping Markdown sanitization as requested: %v", skipSanitize)
		//data := blackfriday.Run(b)
		ghfm := github_flavored_markdown.Markdown(b)
		return ghfm, nil
	}

	log.Printf("DEBUG: Performing Markdown sanitization as requested: %v", !skipSanitize)
	//unsafe := blackfriday.Run(b)
	///data := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	ghfm := github_flavored_markdown.Markdown(b)
	return ghfm, nil

}

// renderDefaultIndexPage is called if the default or user-requested Markdown
// file cannot be opened (e.g., this application binary is being run from
// outside of the directory containing the file)
func renderDefaultIndexPage() string {

	// FIXME: Direct constant access
	return fmt.Sprintf(
		"%s\n%s\n%s",
		htmlHeader,
		htmlFallbackIndexPage,
		htmlFooter,
	)

}

func frontPageHandler(requestedFile string, fallbackContent string, skipSanitize bool) http.HandlerFunc {

	log.Printf("DEBUG: requested file outside of return: %q\n", requestedFile)

	// return "type" of http.HandlerFunc as expected by http.HandleFunc() this
	// function receives `w` and `r` from http.HandleFunc; we do not have to
	// write frontPageHandler() so that it directly receives those `w` and `r`
	// as arguments.
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("DEBUG: frontPageHandler endpoint hit for path: %q\n", r.URL.Path)
		log.Printf("DEBUG: requested file inside return: %q\n", requestedFile)
		//fmt.Fprintf(w, "frontPageHandler endpoint hit")

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
		if r.URL.Path == "/favicon.ico" {
			log.Printf("DEBUG: rejecting request for %q\n", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		filename := requestedFile
		markdownInput, err := loadMarkdown(filename)
		if err != nil {
			log.Printf("Failed to load Markdown file %q: %s", filename, err)
			log.Println("Falling back to static, hard-coded index page.")

			htmlOutput := renderDefaultIndexPage()
			fmt.Fprintf(w, htmlOutput)
			return
		}

		log.Printf("DEBUG: Successfully loaded Markdown file: %q", filename)
		log.Println("DEBUG: Attempting to generate HTML output from loaded Markdown file")

		bytes, err := processMarkdown(markdownInput, skipSanitize)
		if err != nil {
			log.Printf("Failed to parse Markdown file %q: %s", filename, err)
			log.Println("Falling back to static, hard-coded index page.")

			htmlOutput := renderDefaultIndexPage()
			fmt.Fprintf(w, htmlOutput)
			return
		}

		log.Println("DEBUG: Successfully converted HTML from Markdown file")

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

	// SETUP ROUTES

	// Direct request for root of site OR unspecified route (e.g.,"catch-all")
	http.HandleFunc("/", frontPageHandler(readme, htmlFallbackIndexPage, appConfig.SkipMarkdownSanitization))

	// Direct request for readme file
	http.HandleFunc(readme, frontPageHandler(readme, htmlFallbackIndexPage, appConfig.SkipMarkdownSanitization))

	// Direct request for changelog file
	http.HandleFunc(changelog, frontPageHandler(changelog, htmlFallbackIndexPage, appConfig.SkipMarkdownSanitization))

	// TODO: Add useful endpoints for testing here

	// listen on specified port on ALL IP Addresses, block until app is terminated
	listenAddress := fmt.Sprintf(":%d", appConfig.LocalTCPPort)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
