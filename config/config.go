// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package config provides types and functions to collect, validate and apply
// user-provided settings.
package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Updated via Makefile builds. Setting placeholder value here so that
// something resembling a version string will be provided for non-Makefile
// builds.
var version string = "x.y.z"

const myAppName string = "bounce"
const myAppURL string = "https://github.com/atc0005/bounce"

// Default flag settings if not overridden by user input
const (
	defaultLocalTCPPort int = 8000
)

// TCP port ranges
// http://www.iana.org/assignments/port-numbers
// Port numbers are assigned in various ways, based on three ranges: System
// Ports (0-1023), User Ports (1024-49151), and the Dynamic and/or Private
// Ports (49152-65535)
const (
	tcpReservedPort            int = 0
	tcpSystemPortStart         int = 1
	tcpSystemPortEnd           int = 1023
	tcpUserPortStart           int = 1024
	tcpUserPortEnd             int = 49151
	tcpDynamicPrivatePortStart int = 49152
	tcpDynamicPrivatePortEnd   int = 65535
)

// Branding is responsible for emitting application name, version and origin
func Branding() {
	fmt.Fprintf(flag.CommandLine.Output(), "\n%s %s\n%s\n\n", myAppName, version, myAppURL)
}

// Usage is a custom override for the default Help text provided by
// the flag package. Here we prepend some additional metadata to the existing
// output.
func Usage(flagSet *flag.FlagSet) func() {

	return func() {

		Branding()

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"%s\":\n",
			flagSet.Name(),
		)
		flagSet.PrintDefaults()

	}
}

// Config represents the application configuration as specified via
// command-line flags
type Config struct {

	// LocalTCPPort is the TCP port that this application should listen on for
	// incoming requests
	LocalTCPPort int
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"LocalTCPPort: %d",
		c.LocalTCPPort,
	)
}

// NewConfig is a factory function that produces a new Config object based
// on user provided flag values.
func NewConfig() (*Config, error) {

	config := Config{}

	mainFlagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	mainFlagSet.IntVar(
		&config.LocalTCPPort,
		"port",
		defaultLocalTCPPort,
		"TCP port that this application should listen on for incoming HTTP requests.",
	)

	mainFlagSet.Usage = Usage(mainFlagSet)
	// FIXME: Is this needed for any reason since our mainFlagSet has already
	// been parsed?
	//flag.CommandLine = mainFlagSet
	//flag.Parse()
	if err := mainFlagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	// If no errors were encountered during parsing, proceed to validation of
	// configuration settings (both user-specified and defaults)
	if err := validate(config); err != nil {
		return nil, err
	}

	return &config, nil

}

// validate confirms that all config struct fields have reasonable values
func validate(c Config) error {

	switch {

	// WARNING: User opted to use a privileged system port
	case (c.LocalTCPPort >= tcpSystemPortStart) && (c.LocalTCPPort <= tcpSystemPortEnd):

		// DEBUG
		log.Printf(
			"DEBUG: unprivileged system port %d chosen. ports between %d and %d require elevated privileges",
			c.LocalTCPPort,
			tcpSystemPortStart,
			tcpSystemPortEnd,
		)

		// log at WARNING level
		log.Printf(
			"WARNING: Binding to a port < %d requires elevated permissions. If you encounter errors with this application, please re-run this application and specify a port number between %d and %d",
			tcpUserPortStart,
			tcpUserPortStart,
			tcpUserPortEnd,
		)

	// OK: User opted to use a valid and non-privileged port number
	case (c.LocalTCPPort >= tcpUserPortStart) && (c.LocalTCPPort <= tcpUserPortEnd):
		log.Printf(
			"DEBUG: Valid, non-privileged user port between %d and %d configured: %d",
			tcpUserPortStart,
			tcpUserPortEnd,
			c.LocalTCPPort,
		)

	// WARNING: User opted to use a dynamic or private TCP port
	case (c.LocalTCPPort >= tcpDynamicPrivatePortStart) && (c.LocalTCPPort <= tcpDynamicPrivatePortEnd):
		log.Printf(
			"WARNING: Valid, non-privileged, but dynamic/private port between %d and %d configured. This range is reserved for dynamic (usually outgoing) connections. If you encounter errors with this application, please re-run this application and specify a port number between %d and %d",
			tcpUserPortStart,
			tcpUserPortEnd,
			tcpDynamicPrivatePortStart,
			tcpDynamicPrivatePortEnd,
		)

	default:
		log.Printf("invalid port %d specified", c.LocalTCPPort)
		return fmt.Errorf(
			"port %d is not a valid TCP port for this application",
			c.LocalTCPPort,
		)
	}

	// if we made it this far then we signal all is well
	return nil

}
