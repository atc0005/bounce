package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"
)

// NotifyResult wraps the results of goroutine operations to make it easier to
// inspect the status of various tasks so that we can take action on either
// error or success conditions
type NotifyResult struct {
	Val string
	Err error
}

// teamsNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send Microsoft
// Teams messages.
func teamsNotifier(
	ctx context.Context,
	webhookURL string,
	sendTimeout time.Duration,
	retries int,
	retriesDelay int,
	incoming <-chan echoHandlerResponse,
	notifyMgrResultQueue chan<- NotifyResult,
) {

	log.Debug("teamsNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	for {

		// Block while waiting on input
		responseDetails := <-incoming

		log.Debugf("teamsNotifier: Request received: %#v", responseDetails)

		// Wait for specified amount of time before attempting notification.
		// This is done in an effort to prevent unintentional abuse of
		// remote services
		time.Sleep(config.NotifyMgrTeamsNotificationDelay)

		// launch task in separate goroutine
		log.Debug("teamsNotifier: Launching message creation/submission in separate goroutine")
		go func(ctx context.Context, webhookURL string, responseDetails echoHandlerResponse, resultQueue chan<- NotifyResult) {
			ourMessage := createMessage(responseDetails)
			result := NotifyResult{}
			if err := sendMessage(webhookURL, ourMessage, retries, retriesDelay); err != nil {

				result = NotifyResult{
					Err: fmt.Errorf("teamsNotifier: error occurred while trying to send message to Microsoft Teams: %w", err),
				}

				resultQueue <- result
			}

			// Success
			result.Val = "teamsNotifier: Successfully sent message via Microsoft Teams"
			log.Info(result.Val)
			resultQueue <- result
		}(ctx, webhookURL, responseDetails, ourResultQueue)

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine

			result := NotifyResult{
				Err: fmt.Errorf("teamsNotifier: Received Done signal from context"),
			}
			log.Debug(result.Err.Error())
			notifyMgrResultQueue <- result
			return

		case <-time.After(sendTimeout):

			result := NotifyResult{
				Err: fmt.Errorf("teamsNotifier: Timeout reached after %v for sending Microsoft Teams notification", sendTimeout),
			}
			log.Debug(result.Err.Error())
			notifyMgrResultQueue <- result

		case result := <-ourResultQueue:

			if result.Err != nil {
				log.Errorf("teamsNotifier: Error received from ourResultQueue: %v", result.Err)
			} else {
				log.Debugf("teamsNotifier: OK: non-error status received on ourResultQueue: %v", result.Val)
			}

			notifyMgrResultQueue <- result

		}
	}

}

// emailNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send email
// messages.
//
// FIXME: Once the logic is worked out in teamsNotifier, update this function
// to match it
func emailNotifier(ctx context.Context, sendTimeout time.Duration, incoming <-chan echoHandlerResponse, notifyMgrResultQueue chan<- NotifyResult) {

	log.Debug("emailNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	for {

		// Block while waiting on input
		responseDetails := <-incoming

		log.Debugf("emailNotifier: Request received: %#v", responseDetails)

		// Wait for specified amount of time before attempting notification.
		// This is done in an effort to prevent unintentional abuse of
		// remote services
		time.Sleep(config.NotifyMgrEmailNotificationDelay)

		// launch task in a separate goroutine
		go func(resultQueue chan<- NotifyResult) {
			result := NotifyResult{
				Err: fmt.Errorf("emailNotifier: Sending email is not currently enabled"),
			}
			log.Error(result.Err.Error())
			resultQueue <- result
		}(ourResultQueue)

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine

			result := NotifyResult{
				Err: fmt.Errorf("emailNotifier: Received Done signal from context"),
			}
			log.Debug(result.Err.Error())
			notifyMgrResultQueue <- result
			return

		case <-time.After(sendTimeout):

			result := NotifyResult{
				Err: fmt.Errorf("emailNotifier: Timeout reached after %v for sending email notification", sendTimeout),
			}
			log.Debug(result.Err.Error())
			notifyMgrResultQueue <- result

		case result := <-ourResultQueue:

			if result.Err != nil {
				log.Errorf("emailNotifier: Error received from ourResultQueue: %v", result.Err)
			} else {
				log.Debugf("emailNotifier: OK: non-error status received on ourResultQueue: %v", result.Val)
			}

			notifyMgrResultQueue <- result

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
	teamsNotifyResultQueue := make(chan NotifyResult)

	emailNotifyWorkQueue := make(chan echoHandlerResponse)
	emailNotifyResultQueue := make(chan NotifyResult)

	if !cfg.NotifyTeams() && !cfg.NotifyEmail() {
		log.Debug("StartNotifyMgr: Teams and email notifications not requested, not starting notifier goroutines")
	}

	// If enabled, start persistent goroutine to process request details and
	// submit messages to Microsoft Teams.
	if cfg.NotifyTeams() {
		log.Debug("StartNotifyMgr: Starting up teamsNotifier")
		go teamsNotifier(
			ctx,
			cfg.WebhookURL,
			config.NotifyMgrTeamsTimeout,
			cfg.Retries,
			cfg.RetriesDelay,
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

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine
			log.Debug("StartNotifyMgr: Received Done signal from context")
			return

		case result := <-teamsNotifyResultQueue:
			if result.Err != nil {
				log.Errorf("StartNotifyMgr: Error received from teamsNotifyResultQueue: %v", result.Err)
				continue
			}

			log.Debugf("StartNotifyMgr: OK: non-error status received on teamsNotifyResultQueue: %v", result.Val)

		case result := <-emailNotifyResultQueue:
			if result.Err != nil {
				log.Errorf("StartNotifyMgr: Error received from emailNotifyResultQueue: %v", result.Err)
				continue
			}

			log.Debugf("StartNotifyMgr: non-error status received on teamsNotifyResultQueue: %v", result.Val)

		case responseDetails := <-notifyWorkQueue:

			log.Debug("StartNotifyMgr: Input received from notifyWorkQueue")

			// If we don't have *any* notifications enabled we will just pull
			// the incoming item from the the channel and discard it
			if !cfg.NotifyEmail() && !cfg.NotifyTeams() {
				log.Debug("StartNotifyMgr: Notifications are not currently enabled; ignoring notification request")
				continue
			}

			if cfg.NotifyTeams() {
				log.Debug("StartNotifyMgr: Handing off responseDetails to teamsNotifyWorkQueue")
				go func() {
					teamsNotifyWorkQueue <- responseDetails
				}()
			}

			if cfg.NotifyEmail() {
				log.Debug("StartNotifyMgr: Handing off responseDetails to emailNotifyResultQueue")
				go func() {
					emailNotifyWorkQueue <- responseDetails
				}()
			}

			// default:
			// 	log.Debug("StartNotifyMgr: default case statement triggered")
		}

	}
}
