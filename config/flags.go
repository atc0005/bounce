// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"flag"
	"os"
)

// handleFlagsConfig wraps flag setup code into a bundle for potential ease of
// use and future testability
func (c *Config) handleFlagsConfig() error {

	mainFlagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	mainFlagSet.IntVar(
		&c.LocalTCPPort,
		"port",
		defaultLocalTCPPort,
		"TCP port that this application should listen on for incoming HTTP requests.",
	)

	mainFlagSet.StringVar(
		&c.LocalIPAddress,
		"ipaddr",
		defaultLocalIP,
		"Local IP Address that this application should listen on for incoming HTTP requests.",
	)

	mainFlagSet.BoolVar(
		&c.ColorizedJSON,
		"color",
		defaultColorizedJSON,
		"Whether JSON output should be colorized.",
	)

	mainFlagSet.IntVar(
		&c.ColorizedJSONIndent,
		"indent-lvl",
		defaultColorizedJSONIntent,
		"Number of spaces to use when indenting colorized JSON output. Has no effect unless colorized JSON mode is enabled.",
	)

	mainFlagSet.StringVar(
		&c.LogLevel,
		"log-lvl",
		defaultLogLevel,
		"Log message priority filter. Log messages with a lower level are ignored.",
	)

	mainFlagSet.StringVar(
		&c.LogOutput,
		"log-out",
		defaultLogOutput,
		"Log messages are written to this output target",
	)

	mainFlagSet.StringVar(
		&c.LogFormat,
		"log-fmt",
		defaultLogFormat,
		"Log messages are written in this format",
	)

	mainFlagSet.StringVar(
		&c.WebhookURL,
		"webhook-url",
		defaultWebhookURL,
		"The Webhook URL provided by a preconfigured Connector. If specified, this application will attempt to send client request details to the Microsoft Teams channel associated with the webhook URL.",
	)

	mainFlagSet.IntVar(
		&c.Retries,
		"retries",
		defaultRetries,
		"The number of attempts that this application will make to deliver messages before giving up.",
	)

	mainFlagSet.IntVar(
		&c.RetriesDelay,
		"retries-delay",
		defaultRetriesDelay,
		"The number of seconds that this application will wait before making another delivery attempt.",
	)

	mainFlagSet.Usage = Usage(mainFlagSet)

	// TODO: Safe to do this?
	flag.Usage = Usage(mainFlagSet)

	// FIXME: Is this needed for any reason since our mainFlagSet has already
	// been parsed?
	//flag.CommandLine = mainFlagSet
	//flag.Parse()
	if err := mainFlagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil

}
