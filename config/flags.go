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

	mainFlagSet.IntVar(&c.LocalTCPPort, "port", defaultLocalTCPPort, portFlagHelp)
	mainFlagSet.StringVar(&c.LocalIPAddress, "ipaddr", defaultLocalIP, localIPAddressFlagHelp)
	mainFlagSet.BoolVar(&c.ColorizedJSON, "color", defaultColorizedJSON, colorizedJSONFlagHelp)
	mainFlagSet.IntVar(&c.ColorizedJSONIndent, "indent-lvl", defaultColorizedJSONIntent, colorizedJSONIndentFlagHelp)
	mainFlagSet.StringVar(&c.LogLevel, "log-lvl", defaultLogLevel, logLevelFlagHelp)
	mainFlagSet.StringVar(&c.LogOutput, "log-out", defaultLogOutput, logOutputFlagHelp)
	mainFlagSet.StringVar(&c.LogFormat, "log-fmt", defaultLogFormat, logFormatFlagHelp)
	mainFlagSet.StringVar(&c.WebhookURL, "webhook-url", defaultWebhookURL, webhookURLFlagHelp)
	mainFlagSet.IntVar(&c.Retries, "retries", defaultRetries, retriesFlagHelp)
	mainFlagSet.IntVar(&c.RetriesDelay, "retries-delay", defaultRetriesDelay, retriesDelayFlagHelp)

	mainFlagSet.Usage = Usage(mainFlagSet)

	// TODO: Safe to do this?
	flag.Usage = Usage(mainFlagSet)

	// FIXME: Is this needed for any reason since our mainFlagSet has already
	// been parsed?
	// flag.CommandLine = mainFlagSet
	// flag.Parse()
	if err := mainFlagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil

}
