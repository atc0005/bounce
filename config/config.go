// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bridge
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package config provides types and functions to collect, validate and apply
// user-provided settings.
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Updated via Makefile builds. Setting placeholder value here so that
// something resembling a version string will be provided for non-Makefile
// builds.
var version string = "x.y.z"

const myAppName string = "bounce"
const myAppURL string = "https://github.com/atc0005/bounce"

const defaultLocalTCPPort int = 8000
const defaultInputMarkdownFile string = "README.md"
const defaultSkipMarkdownSanitization bool = false

// Branding is responsible for emitting application name, version and origin
func Branding() {
	fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s\n%s\n\n", myAppName, version, myAppURL)
}

// Usage is a custom override for the default Help text provided by
// the flag package. Here we prepend some additional metadata to the existing
// output.
func Usage(flagSet *flag.FlagSet) func() {

	return func() {

		myBinaryName := filepath.Base(os.Args[0])

		Branding()

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"%s %s\":\n",
			myBinaryName,
			flagSet.Name(),
		)
		flagSet.PrintDefaults()

	}
}

// Config represents the application configuration as specified via
// command-line flags
type Config struct {

	// InputFile represents the full path to an input Markdown file. This is
	// usually the path to this repo's README.md file.
	InputFile string

	// SkipMarkdownSanitization indicates whether sanitization of Markdown
	// input should be skipped. The default is to perform this sanitization to
	// help protect against untrusted input.
	SkipMarkdownSanitization bool

	// LocalTCPPort is the TCP port that this application should listen on for
	// incoming requests
	LocalTCPPort int
}

// NewConfig is a factory function that produces a new Config object based
// on user provided flag values.
func NewConfig() (*Config, error) {

	config := Config{}

	mainFlagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	mainFlagSet.Usage = Usage(mainFlagSet)
	flag.CommandLine = mainFlagSet
	flag.Parse()

	mainFlagSet.StringVar(&config.InputFile, "input-file", defaultInputMarkdownFile, "Path to Markdown file to process and display. The default is this repo's README.md file.")
	mainFlagSet.IntVar(&config.LocalTCPPort, "port", defaultLocalTCPPort, "Number of files of the same file size needed before duplicate validation logic is applied.")
	mainFlagSet.BoolVar(&config.SkipMarkdownSanitization, "skip-sanitize", defaultSkipMarkdownSanitization, "Perform recursive search into subdirectories per provided path.")

	return &config, nil

}
