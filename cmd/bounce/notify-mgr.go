package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"
)

// teamsNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send Microsoft
// Teams messages.
func teamsNotifier(ctx context.Context, webhookURL string, sendTimeout time.Duration, incoming <-chan echoHandlerResponse, notifyMgrResultQueue chan<- error) {

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan error)

	for {

		// Block while waiting on input
		responseDetails := <-incoming

		log.Debugf("teamsNotifier: Request received: %#v", responseDetails)

		// launch task in separate goroutine
		go func(ctx context.Context, webhookURL string, responseDetails echoHandlerResponse, result chan<- error) {
			ourMessage := createMessage(responseDetails)
			if err := sendMessage(webhookURL, ourMessage); err != nil {
				result <- fmt.Errorf("teamsNotifier: error occurred while trying to send message to Microsoft Teams: %w", err)
			}

			// Success
			log.Info("teamsNotifier: Successfully sent message via Microsoft Teams")
			result <- nil
		}(ctx, webhookURL, responseDetails, ourResultQueue)

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			// log.Debug("teamsNotifier: Received Done signal from context")
			notifyMgrResultQueue <- fmt.Errorf("teamsNotifier: Received Done signal from context")
			return

		case <-time.After(sendTimeout):

			// log.Debugf("teamsNotifier: Timeout reached after %v for sending Microsoft Teams notification", sendTimeout)
			notifyMgrResultQueue <- fmt.Errorf("teamsNotifier: Timeout reached after %v for sending Microsoft Teams notification", sendTimeout)

		case err := <-ourResultQueue:
			if err != nil {
				// log.Errorf("teamsNotifier: Error received from ourResultQueue: %v", err.Error())
				notifyMgrResultQueue <- fmt.Errorf("teamsNotifier: Error received from ourResultQueue: %v", err.Error())
			}

			// log.Debug("teamsNotifier: non-error status received on ourResultQueue")
			notifyMgrResultQueue <- fmt.Errorf("teamsNotifier: non-error status received on ourResultQueue")

		}
	}

}

// emailNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send email
// messages.
//
// FIXME: Once the logic is worked out in teamsNotifier, update this function
// to match it
func emailNotifier(ctx context.Context, sendTimeout time.Duration, incoming <-chan echoHandlerResponse, notifyMgrResultQueue chan<- error) {

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan error)

	for {

		// Block while waiting on input
		responseDetails := <-incoming

		log.Debugf("emailNotifier: Request received: %#v", responseDetails)

		// launch task in a separate goroutine
		go func() {
			// log.Error("emailNotifier: Sending email is not currently enabled.")
			notifyMgrResultQueue <- fmt.Errorf("emailNotifier: Sending email is not currently enabled.")
		}()

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			// log.Debug("emailNotifier: Received Done signal from context")
			notifyMgrResultQueue <- fmt.Errorf("emailNotifier: Received Done signal from context")
			return

		case <-time.After(sendTimeout):

			// log.Debugf("emailNotifier: Timeout reached after %v for sending Microsoft Teams notification", sendTimeout)
			notifyMgrResultQueue <- fmt.Errorf("emailNotifier: Timeout reached after %v for sending Microsoft Teams notification", sendTimeout)

		case err := <-ourResultQueue:
			if err != nil {
				// log.Errorf("emailNotifier: Error received from ourResultQueue: %v", err.Error())
				notifyMgrResultQueue <- fmt.Errorf("emailNotifier: Error received from ourResultQueue: %v", err.Error())
			}

			// log.Debug("emailNotifier: non-error status received on ourResultQueue")
			notifyMgrResultQueue <- fmt.Errorf("emailNotifier: non-error status received on ourResultQueue")

		}
	}

}

// StartNotifyMgr receives echoHandlerResponse values from a receive-only
// incoming queue of echoHandlerResponse values and sends notifications to any
// enabled service (e.g., Microsoft Teams).
func StartNotifyMgr(ctx context.Context, cfg *config.Config, notifyWorkQueue <-chan echoHandlerResponse) {

	// Create channels to hand-off echoHandlerResponse values for
	// processing. Due to my ignorance of channels, I believe that I'll need
	// separate channels for each service. E.g., one channel for Microsoft
	// Teams outgoing notifications, another for email and so on.

	teamsNotifyWorkQueue := make(chan echoHandlerResponse)
	teamsNotifyResultQueue := make(chan error)

	emailNotifyWorkQueue := make(chan echoHandlerResponse)
	emailNotifyResultQueue := make(chan error)

	// If enabled, start persistent goroutine to process request details and
	// submit messages to Microsoft Teams.
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

	// If enabled, start persistent goroutine to process request details and
	// submit messages by email.
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

		// If we don't have *any* notifications enabled we will just pull
		// the incoming item from the the channel and discard it
		if !cfg.NotifyEmail() && !cfg.NotifyTeams() {
			log.Debug("StartNotifyMgr: Notifications are not currently enabled; ignoring notification request")
			continue
		}

		if cfg.NotifyTeams() {
			go func() {
				log.Debug("StartNotifyMgr: Handed off responseDetails to teamsNotifyWorkQueue")
				teamsNotifyWorkQueue <- responseDetails
			}()
		}

		if cfg.NotifyEmail() {
			go func() {
				log.Debug("StartNotifyMgr: Handed off responseDetails to emailNotifyResultQueue")
				emailNotifyWorkQueue <- responseDetails
			}()
		}

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			log.Debug("StartNotifyMgr: Received Done signal from context")
			return

		case err := <-teamsNotifyResultQueue:
			if err != nil {
				log.Errorf("StartNotifyMgr: Error received from teamsNotifyResultQueue: %v", err.Error())
				continue
			}

			log.Debug("StartNotifyMgr: non-error status received on teamsNotifyResultQueue")

		case err := <-emailNotifyResultQueue:
			if err != nil {
				log.Errorf("StartNotifyMgr: Error received from emailNotifyResultQueue: %v", err.Error())
				continue
			}

			log.Debug("StartNotifyMgr: non-error status received on teamsNotifyResultQueue")

		default:
			log.Debug("StartNotifyMgr: default case statement triggered")
		}

	}
}
