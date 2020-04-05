package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"
)

// teamsNotifier handles generating Microsoft Teams notifications from
// incoming client request details.
func teamsNotifier(ctx context.Context, webhookURL string, sendTimeout time.Duration, incoming <-chan echoHandlerResponse, result chan<- error) {

	for {

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			log.Debug("teamsNotifier: Received Done signal from context")
			return

		// Block waiting on input from notifyWorkQueue channel
		case responseDetails := <-incoming:

			log.Debugf("teamsNotifier: teamsNotifierRequest received by teams notification goroutine: %#v", responseDetails)

			ourMessage := createMessage(responseDetails)
			if err := sendMessage(webhookURL, ourMessage); err != nil {
				result <- fmt.Errorf("teamsNotifier: error occurred while trying to send message to Microsoft Teams: %w", err)
			}

			// Success
			log.Info("teamsNotifier: Successfully sent message via Microsoft Teams")
			result <- nil

		case <-time.After(sendTimeout):

			log.Info("teamsNotifier: Timeout reached after for sending Microsoft Teams notification")

			// FIXME
			// 	Q: send message back to log?
			//	A: is there a problem with using logging here within a goroutine?
			// FIXME
			//	Q: should a return be used here?
			// 	A: presumably not as it would cancel the reception of further messages for processing

		}
	}

}

// spin off goroutine to create and send email messages
func emailNotifier(ctx context.Context, sendTimeout time.Duration, incoming <-chan echoHandlerResponse, result chan<- error) {

	for {

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			log.Debug("emailNotifier: Received Done signal from context")
			return

		// Block waiting on input from notifyWorkQueue channel
		case responseDetails := <-incoming:
			log.Debugf("emailNotifier: Request received by email notification goroutine: %#v", responseDetails)

			errMsg := "emailNotifier: Sending email is not currently enabled."
			log.Error(errMsg)
			result <- fmt.Errorf(errMsg)

		case <-time.After(sendTimeout):
			log.Debug("emailNotifier: Timeout reached for sending email notification.")

			// FIXME: send message back to log?
			// FIXME: should a return be used here?
		}
	}

}

// StartNotifyMgr receives echoHandlerResponse values from a receive-only
// incoming queue of echoHandlerResponse values and sends notifications to any
// enabled service (e.g., Microsoft Teams).
func StartNotifyMgr(ctx context.Context, cfg *config.Config, notifyWorkQueue <-chan echoHandlerResponse) {

	// https://gobyexample.com/channel-directions
	//
	// func pong(pings <-chan string, pongs chan<- string) {
	// 	msg := <-pings
	// 	pongs <- msg
	// }

	// Create channels to hand-off echoHandlerResponse values for
	// processing. Due to my ignorance of channels, I believe that I'll need
	// separate channels for each service. E.g., one channel for Microsoft
	// Teams outgoing notifications, another for email and so on.

	teamsNotifyWorkQueue := make(chan echoHandlerResponse)
	teamsNotifyResultQueue := make(chan error)

	emailNotifyWorkQueue := make(chan echoHandlerResponse)
	emailNotifyResultQueue := make(chan error)

	// Send request details to Microsoft Teams if webhook URL set
	if cfg.NotifyTeams() {
		log.Debug("StartNotifyMgr: Starting up teamsNotifier")
		go teamsNotifier(
			ctx,
			cfg.WebhookURL,
			config.NotifyMgrTeamsTimeout,
			teamsNotifyWorkQueue,
			teamsNotifyResultQueue,
		)
	}

	if cfg.NotifyEmail() {
		log.Debug("StartNotifyMgr: Starting up emailNotifier")
		go emailNotifier(
			ctx,
			config.NotifyMgrEmailTimeout,
			emailNotifyWorkQueue,
			emailNotifyResultQueue,
		)
	}

	for {

		// Block waiting on input from notifyWorkQueue channel
		log.Debug("StartNotifyMgr: Waiting on input from notifyWorkQueue")
		responseDetails := <-notifyWorkQueue
		log.Debug("StartNotifyMgr: Input received from notifyWorkQueue")

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			log.Debug("StartNotifyMgr: Received Done signal from context")
			return

		// Attempt to send response details to goroutine responsible for
		// generating Microsoft Teams messages
		case teamsNotifyWorkQueue <- responseDetails:
			log.Debug("StartNotifyMgr: Handed off responseDetails to teamsNotifyWorkQueue")
			// success; now what?

		case err := <-teamsNotifyResultQueue:
			// do something based on success or failure sending to Teams
			if err != nil {
				log.Errorf("StartNotifyMgr: Error received from teamsNotifyResultQueue: %v", err.Error())
			}

		// Attempt to send response details to goroutine responsible for
		// generating email messages
		case emailNotifyWorkQueue <- responseDetails:
			log.Debug("StartNotifyMgr: Handed off responseDetails to emailNotifyWorkQueue")
			// success; now what?

		case err := <-emailNotifyResultQueue:
			// do something based on success or failure sending to Teams
			if err != nil {
				log.Errorf("StartNotifyMgr: Error received from emailNotifyResultQueue: %v", err.Error())
			}
		}

	}
}
