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
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/discard"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/text"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify"
	send2teams "github.com/atc0005/send2teams/teams"
)

// version is updated via Makefile builds by referencing the fully-qualified
// path to this variable, including the package. We set a placeholder value so
// that something resembling a version string will be provided for
// non-Makefile builds.
var version string = "x.y.z"

// MyAppName is the public name of this application
const MyAppName string = "bounce"

// MyAppURL is the location of the repo for this application
const MyAppURL string = "https://github.com/atc0005/bounce"

const (
	portFlagHelp                = "TCP port that this application should listen on for incoming HTTP requests."
	localIPAddressFlagHelp      = "Local IP Address that this application should listen on for incoming HTTP requests."
	colorizedJSONFlagHelp       = "Whether JSON output should be colorized."
	colorizedJSONIndentFlagHelp = "Number of spaces to use when indenting colorized JSON output. Has no effect unless colorized JSON mode is enabled."
	logLevelFlagHelp            = "Log message priority filter. Log messages with a lower level are ignored."
	logOutputFlagHelp           = "Log messages are written to this output target"
	logFormatFlagHelp           = "Log messages are written in this format"
	webhookURLFlagHelp          = "The Webhook URL provided by a preconfigured Connector. If specified, this application will attempt to send client request details to the Microsoft Teams channel associated with the webhook URL."
	retriesFlagHelp             = "The number of attempts that this application will make to deliver messages before giving up."
	retriesDelayFlagHelp        = "The number of seconds that this application will wait before making another delivery attempt."
)

// Default flag settings if not overridden by user input
const (
	defaultLocalTCPPort        int    = 8000
	defaultLocalIP             string = "localhost"
	defaultColorizedJSON       bool   = false
	defaultColorizedJSONIntent int    = 2
	defaultLogLevel            string = "info"
	defaultLogOutput           string = "stdout"
	defaultLogFormat           string = "text"
	defaultWebhookURL          string = ""
	defaultRetries             int    = 2
	defaultRetriesDelay        int    = 2
)

// Timeout settings applied to our instance of http.Server
const (
	HTTPServerReadHeaderTimeout time.Duration = 20 * time.Second
	HTTPServerReadTimeout       time.Duration = 1 * time.Minute
	HTTPServerWriteTimeout      time.Duration = 2 * time.Minute
)

// HTTPServerShutdownTimeout is used by the graceful shutdown process to
// control how long the shutdown process should wait before forcefully
// terminating.
const HTTPServerShutdownTimeout time.Duration = 30 * time.Second

// NotifyMgrServicesShutdownTimeout is used by the NotifyMgr to determine how
// long it should wait for results from each notifier or notifier "service"
// before continuing on with the shutdown process.
const NotifyMgrServicesShutdownTimeout time.Duration = 2 * time.Second

// Timing-related settings (delays, timeouts) used by our notification manager
// when using goroutines to concurrently process notification requests.
const (

	// NotifyMgrTeamsTimeout is the timeout setting applied to each Microsoft
	// Teams notification attempt. This value does NOT take into account the
	// number of configured retries and retry delays. The final value timeout
	// applied to each notification attempt should be based on those
	// calculations. The TeamsTimeout method does just that.
	NotifyMgrTeamsTimeout time.Duration = 10 * time.Second

	// NotifyMgrTeamsSendAttemptTimeout

	// NotifyMgrEmailTimeout is the timeout setting applied to each email
	// notification attempt.
	// TODO: Email support is not (as of this writing) available. This is a
	// stub entry to satisfy stub functionality for later use.
	NotifyMgrEmailTimeout time.Duration = 30 * time.Second

	// NotifyStatsMonitorDelay limits notification stats logging to no more
	// often than this duration. This limiter is to keep from logging the
	// details so often that the information simply becomes noise.
	NotifyStatsMonitorDelay time.Duration = 30 * time.Second

	// NotifyQueueMonitorDelay limits notification queue stats logging to no
	// more often than this duration. This limiter is to keep from logging the
	// details so often that the information simply becomes noise.
	NotifyQueueMonitorDelay time.Duration = 15 * time.Second

	// NotifyMgrTeamsNotificationDelay is the delay between Microsoft Teams
	// notification attempts. This delay is intended to help prevent
	// unintentional abuse of remote services.
	NotifyMgrTeamsNotificationDelay time.Duration = 5 * time.Second

	// NotifyMgrEmailNotificationDelay is the delay between email notification
	// attempts. This delay is intended to help prevent unintentional abuse of
	// remote services.
	NotifyMgrEmailNotificationDelay time.Duration = 5 * time.Second
)

// NotifyMgrQueueDepth is the number of items allowed into the queue/channel
// at one time. Senders with items for the notification "pipeline" that do not
// fit within the allocated space will block until space in the queue opens.
// Best practice for channels advocates that a smaller number is better than a
// larger one, so YMMV if this is set either too high or too low.
//
// Brief testing (as of this writing) shows that a depth as low as 1 works for
// our purposes, but results in a greater number of stalled goroutines waiting
// to place items into the queue.
const NotifyMgrQueueDepth int = 5

// ReadHeaderTimeout:

// TCP port ranges
// http://www.iana.org/assignments/port-numbers
// Port numbers are assigned in various ways, based on three ranges: System
// Ports (0-1023), User Ports (1024-49151), and the Dynamic and/or Private
// Ports (49152-65535)
const (
	TCPReservedPort            int = 0
	TCPSystemPortStart         int = 1
	TCPSystemPortEnd           int = 1023
	TCPUserPortStart           int = 1024
	TCPUserPortEnd             int = 49151
	TCPDynamicPrivatePortStart int = 49152
	TCPDynamicPrivatePortEnd   int = 65535
)

// Log levels
const (
	// https://godoc.org/github.com/apex/log#Level

	// LogLevelFatal is used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LogLevelFatal string = "fatal"

	// LogLevelError is for errors that should definitely be noted.
	LogLevelError string = "error"

	// LogLevelWarn is for non-critical entries that deserve eyes.
	LogLevelWarn string = "warn"

	// LogLevelInfo is for general application operational entries.
	LogLevelInfo string = "info"

	// LogLevelDebug is for debug-level messages and is usually enabled
	// when debugging. Very verbose logging.
	LogLevelDebug string = "debug"
)

// 	apex/log Handlers
// ---------------------------------------------------------
// cli - human-friendly CLI output
// discard - discards all logs
// es - Elasticsearch handler
// graylog - Graylog handler
// json - JSON output handler
// kinesis - AWS Kinesis handler
// level - level filter handler
// logfmt - logfmt plain-text formatter
// memory - in-memory handler for tests
// multi - fan-out to multiple handlers
// papertrail - Papertrail handler
// text - human-friendly colored output
// delta - outputs the delta between log calls and spinner
const (
	// LogFormatCLI provides human-friendly CLI output
	LogFormatCLI string = "cli"

	// LogFormatJSON provides JSON output
	LogFormatJSON string = "json"

	// LogFormatLogFmt provides logfmt plain-text output
	LogFormatLogFmt string = "logfmt"

	// LogFormatText provides human-friendly colored output
	LogFormatText string = "text"

	// LogFormatDiscard discards all logs
	LogFormatDiscard string = "discard"
)

const (

	// LogOutputStdout represents os.Stdout
	LogOutputStdout string = "stdout"

	// LogOutputStderr represents os.Stderr
	LogOutputStderr string = "stderr"
)

// MessageTrailer generates a branded "footer" for use with notifications.
func MessageTrailer() string {
	return fmt.Sprintf("Message generated by [%s](%s) (%s)", MyAppName, MyAppURL, version)
}

// Branding is responsible for emitting application name, version and origin
func Branding() string {
	return fmt.Sprintf("\n%s %s\n%s\n\n", MyAppName, version, MyAppURL)
}

// Usage is a custom override for the default Help text provided by
// the flag package. Here we prepend some additional metadata to the existing
// output.
func Usage(flagSet *flag.FlagSet) func() {

	return func() {

		fmt.Fprint(flag.CommandLine.Output(), Branding())

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of \"%s\":\n",
			flagSet.Name(),
		)
		flagSet.PrintDefaults()

	}
}

// Config represents the application configuration as specified via
// command-line flags
type Config struct {

	// Retries is the number of attempts that this application will make
	// to deliver messages before giving up.
	Retries int

	// RetriesDelay is the number of seconds to wait between retry attempts.
	RetriesDelay int

	// LocalTCPPort is the TCP port that this application should listen on for
	// incoming requests
	LocalTCPPort int

	// LocalIPAddress is the IP Address that this application should listen on
	// for incoming requests
	LocalIPAddress string

	// ColorizedJSON indicates whether JSON output should be colorized.
	// Coloring the output could aid in in quick visual evaluation of incoming
	// payloads
	ColorizedJSON bool

	// ColorizedJSONIndent controls how many spaces are used when indenting
	// colorized JSON output. If ColorizedJSON is not enabled, this setting
	// has no effect.
	ColorizedJSONIndent int

	// LogLevel is the chosen logging level
	LogLevel string

	// LogOutput is one of the standard application outputs, stdout or stderr
	// FIXME: Needs better description
	LogOutput string

	// LogFormat controls which output format is used for log messages
	// generated by this application. This value is from a smaller subset
	// of the formats supported by the third-party leveled-logging package
	// used by this application.
	LogFormat string

	// WebhookURL is the full URL used to submit messages to the Teams channel
	// This URL is in the form of https://outlook.office.com/webhook/xxx or
	// https://outlook.office365.com/webhook/xxx. This URL is REQUIRED in
	// order for this application to function and needs to be created in
	// advance by adding/configuring a Webhook Connector in a Microsoft Teams
	// channel that you wish to submit messages to using this application.
	WebhookURL string
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"LocalTCPPort: %d, LocalIPAddress: %s, ColorizedJSON: %t, ColorizedJSONIndent: %d, LogLevel: %s, LogOutput: %s, LogFormat: %s, WebhookURL: %s",
		c.LocalTCPPort,
		c.LocalIPAddress,
		c.ColorizedJSON,
		c.ColorizedJSONIndent,
		c.LogLevel,
		c.LogOutput,
		c.LogFormat,
		c.WebhookURL,
	)
}

// NotifyTeams indicates whether or not notifications should be sent to a
// Microsoft Teams channel.
func (c Config) NotifyTeams() bool {

	// Assumption: config.validate() has already been called for the existing
	// instance of the Config type and this method is now being called by
	// later stages of the codebase to determine only whether an attempt
	// should be made to send a message to Teams.

	// For now, use the same logic that validate() uses to determine whether
	// validation checks should be run: Is c.WebhookURL set to a non-empty
	// string.
	return c.WebhookURL != ""

}

// NotifyEmail indicates whether or not notifications should be generated and
// sent via email to specified recipients.
func (c Config) NotifyEmail() bool {

	// TODO: Add support for email notifications. For now, this method is a
	// placeholder to allow logic for future notification support to be
	// written.
	return false

}

// TeamsTimeout accepts the next scheduled notification, the number of
// Microsoft Teams message submission retries and the delay between each
// attempt and returns the timeout value for the entire message submission
// process, including the initial attempt and all retry attempts.
//
// This overall timeout value is computed using multiple values; (1) the base
// timeout value for a single message submission attempt, (2) the next
// scheduled notification (which was created using the configured delay we
// wish to force between message submission attempts), (3) the total number of
// retries allowed, (4) the delay between retry attempts
func TeamsTimeout(schedule time.Time, retries int, retriesDelay int) time.Duration {

	timeoutValue := (NotifyMgrTeamsTimeout + time.Until(schedule)) +
		(time.Duration(retriesDelay) * time.Duration(retries))

	// Note: This seems to allow the app to make it all the way to and execute
	// goteamsnotify mstClient.SendWithContext() once before the context
	// timeout is triggered and shuts everything down
	// timeoutValue := 6000 * time.Millisecond

	// ... to make it to
	// "sendMessage: Waiting for either context or notificationDelayTimer"
	// before the context expires (0 executions of SendWithContext()).
	// timeoutValue := 5010 * time.Millisecond

	return timeoutValue
}

// NewConfig is a factory function that produces a new Config object based
// on user provided flag values.
func NewConfig() (*Config, error) {

	config := Config{}

	if err := config.handleFlagsConfig(); err != nil {
		return nil, fmt.Errorf("error encountered configuring flags: %w", err)
	}

	// Apply initial logging settings based on any provided CLI flags
	config.configureLogging()

	// If no errors were encountered during parsing, proceed to validation of
	// configuration settings (both user-specified and defaults)
	if err := validate(config); err != nil {
		flag.Usage()
		return nil, err
	}

	return &config, nil

}

// configureLogging is a wrapper function to enable setting requested logging
// settings.
func (c Config) configureLogging() {

	var logOutput *os.File

	switch c.LogOutput {
	case LogOutputStderr:
		logOutput = os.Stderr
	case LogOutputStdout:
		logOutput = os.Stdout
	}

	switch c.LogFormat {
	case LogFormatCLI:
		log.SetHandler(cli.New(logOutput))
	case LogFormatJSON:
		log.SetHandler(json.New(logOutput))
	case LogFormatLogFmt:
		log.SetHandler(logfmt.New(logOutput))
	case LogFormatText:
		log.SetHandler(text.New(logOutput))
	case LogFormatDiscard:
		log.SetHandler(discard.New())
	}

	switch c.LogLevel {
	case LogLevelFatal:
		log.SetLevel(log.FatalLevel)
	case LogLevelError:
		log.SetLevel(log.ErrorLevel)
	case LogLevelWarn:
		log.SetLevel(log.WarnLevel)
	case LogLevelInfo:
		log.SetLevel(log.InfoLevel)
	case LogLevelDebug:
		log.SetLevel(log.DebugLevel)
	}
}

// validate confirms that all config struct fields have reasonable values
func validate(c Config) error {

	switch {

	// WARNING: User opted to use a privileged system port
	case (c.LocalTCPPort >= TCPSystemPortStart) && (c.LocalTCPPort <= TCPSystemPortEnd):

		log.Debugf(
			"unprivileged system port %d chosen. ports between %d and %d require elevated privileges",
			c.LocalTCPPort,
			TCPSystemPortStart,
			TCPSystemPortEnd,
		)

		// log at WARNING level
		log.Warnf(
			"Binding to a port < %d requires elevated permissions. If you encounter errors with this application, please re-run this application and specify a port number between %d and %d",
			TCPUserPortStart,
			TCPUserPortStart,
			TCPUserPortEnd,
		)

	// OK: User opted to use a valid and non-privileged port number
	case (c.LocalTCPPort >= TCPUserPortStart) && (c.LocalTCPPort <= TCPUserPortEnd):
		log.Debugf(
			"Valid, non-privileged user port between %d and %d configured: %d",
			TCPUserPortStart,
			TCPUserPortEnd,
			c.LocalTCPPort,
		)

	// WARNING: User opted to use a dynamic or private TCP port
	case (c.LocalTCPPort >= TCPDynamicPrivatePortStart) && (c.LocalTCPPort <= TCPDynamicPrivatePortEnd):
		log.Warnf(
			"WARNING: Valid, non-privileged, but dynamic/private port between %d and %d configured. This range is reserved for dynamic (usually outgoing) connections. If you encounter errors with this application, please re-run this application and specify a port number between %d and %d",
			TCPUserPortStart,
			TCPUserPortEnd,
			TCPDynamicPrivatePortStart,
			TCPDynamicPrivatePortEnd,
		)

	default:
		log.Debugf("invalid port %d specified", c.LocalTCPPort)
		return fmt.Errorf(
			"port %d is not a valid TCP port for this application",
			c.LocalTCPPort,
		)
	}

	if c.LocalIPAddress == "" {
		return fmt.Errorf("local IP Address not provided")
	}

	// TODO: Consider also throwing an error if this is set without also
	// setting c.ColorizedJSON
	if c.ColorizedJSONIndent <= 0 {
		return fmt.Errorf(
			"invalid indent level chosen for colorized output: %d",
			c.ColorizedJSONIndent,
		)
	}

	switch c.LogLevel {
	case LogLevelFatal:
	case LogLevelError:
	case LogLevelWarn:
	case LogLevelInfo:
	case LogLevelDebug:
	default:
		return fmt.Errorf("invalid option %q provided for log level",
			c.LogLevel)
	}

	switch c.LogOutput {
	case LogOutputStderr:
	case LogOutputStdout:
	default:
		return fmt.Errorf("invalid option %q provided for log output",
			c.LogOutput)
	}

	switch c.LogFormat {
	case LogFormatCLI:
	case LogFormatJSON:
	case LogFormatLogFmt:
	case LogFormatText:
	case LogFormatDiscard:
	default:
		return fmt.Errorf("invalid option %q provided for log format",
			c.LogFormat)
	}

	// LogFormat

	// Not having a webhook URL is a valid choice. Perform validation if value
	// is provided.
	if c.WebhookURL != "" {

		// TODO: Do we really need both of these?
		if ok, err := goteamsnotify.IsValidWebhookURL(c.WebhookURL); !ok {
			return err
		}

		if err := send2teams.ValidateWebhook(c.WebhookURL); err != nil {
			return err
		}

	}

	// if we made it this far then we signal all is well
	return nil

}
