package main

import (
	"context"
	"fmt"
	"runtime"
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

// NotifyQueue represents a channel used to queue input data and responses
// between the main application, the notifications manager and "notifiers".
type NotifyQueue struct {
	// The name of a queue. This is intended for display in log messages.
	Name string

	// Channel is a channel used to transport input data and responses.
	Channel interface{}
}

// notifyQueueMonitor accepts a context and one or many NotifyQueue values to
// monitor for items yet to be processed. notifyQueueMonitor is intended to be
// run as a goroutine

func notifyQueueMonitor(ctx context.Context, delay time.Duration, notifyQueues ...NotifyQueue) {

	if len(notifyQueues) == 0 {
		log.Error("received empty list of notifyQueues to monitor, exiting")
		return
	}

	log.Debug("notifyQueueMonitor: Running")

	t := time.NewTimer(delay)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			// returning not to leak the goroutine
			ctxErr := ctx.Err()
			log.Debugf("notifyQueueMonitor: Received Done signal: %v, shutting down ...", ctxErr.Error())
			return

		// Show stats only for queues with content
		case <-t.C:

			var itemsFound bool
			//log.Debugf("Length of queues: %d", len(queues))
			for _, notifyQueue := range notifyQueues {

				var queueLength int
				switch queue := notifyQueue.Channel.(type) {

				// FIXME: Is there a generic way to match any channel type
				// here in order to calculate the length?
				case chan echoHandlerResponse:
					queueLength = len(queue)

				case <-chan echoHandlerResponse:
					queueLength = len(queue)

				case chan NotifyResult:
					queueLength = len(queue)

				default:
					log.Warn("Default case triggered (this should not happen")
					log.Warnf("Name of channel: %s", notifyQueue.Name)

				}

				if queueLength > 0 {
					itemsFound = true
					log.Warnf("notifyQueueMonitor: %d items in %s",
						queueLength, notifyQueue.Name)
					log.Warnf("notifyQueueMonitor: %d goroutines running",
						runtime.NumGoroutine())
					continue
				}

			}

			if !itemsFound {
				log.Warn("notifyQueueMonitor: 0 items in any monitored queues")
				log.Warnf("notifyQueueMonitor: %d goroutines running",
					runtime.NumGoroutine())
			}
		}
	}

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
	done chan<- struct{},
) {

	log.Debug("teamsNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	for {

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine

			ctxErr := ctx.Err()
			result := NotifyResult{
				Val: fmt.Sprintf("teamsNotifier: Received Done signal: %v, shutting down", ctxErr.Error()),
			}
			log.Debug(result.Val)

			log.Debug("teamsNotifier: Sending back results")
			notifyMgrResultQueue <- result

			log.Debug("teamsNotifier: Closing notifyMgrResultQueue channel to signal shutdown")
			close(notifyMgrResultQueue)

			log.Debug("teamsNotifier: Closing done channel to signal shutdown")
			close(done)
			log.Debug("teamsNotifier: done channel closed, returning")
			return

		case responseDetails := <-incoming:

			log.Debugf("teamsNotifier: Request received: %#v", responseDetails)

			// Wait for specified amount of time before attempting notification.
			// This is done in an effort to prevent unintentional abuse of
			// remote services
			log.Debugf("teamsNotifier: Waiting for %v before processing new request", config.NotifyMgrTeamsNotificationDelay)
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
				result.Val = "teamsNotifier: Successfully sent message to Microsoft Teams"
				log.Info(result.Val)
				resultQueue <- result
			}(ctx, webhookURL, responseDetails, ourResultQueue)

			// Wait for either the timeout to occur OR a result to come back
			// from the attempt to send a Teams message.

			t := time.NewTimer(sendTimeout)
			defer t.Stop()

			select {
			case <-t.C:

				result := NotifyResult{
					Err: fmt.Errorf("teamsNotifier: Timeout reached after %v for sending Microsoft Teams notification", sendTimeout),
				}
				log.Debug(result.Err.Error())
				notifyMgrResultQueue <- result

				// TODO
				// Q: How to actually abandon the Teams message submission?
				// A: Pass context on to sendMessage() function?
				//    Update that function to use context?
				//    Call cancel() and then use continue to loop back around?

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

}

// emailNotifier is a persistent goroutine used to receive incoming
// notification requests and spin off goroutines to create and send email
// messages.
//
// FIXME: Once the logic is worked out in teamsNotifier, update this function
// to match it
func emailNotifier(ctx context.Context, sendTimeout time.Duration, incoming <-chan echoHandlerResponse, notifyMgrResultQueue chan<- NotifyResult, done chan<- struct{}) {

	log.Debug("emailNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	for {

		select {

		case <-ctx.Done():
			// returning not to leak the goroutine

			ctxErr := ctx.Err()
			result := NotifyResult{
				Val: fmt.Sprintf("emailNotifier: Received Done signal: %v, shutting down", ctxErr.Error()),
			}
			log.Debug(result.Val)

			log.Debug("emailNotifier: Sending back results")
			notifyMgrResultQueue <- result

			log.Debug("emailNotifier: Closing notifyMgrResultQueue channel to signal shutdown")
			close(notifyMgrResultQueue)

			log.Debug("emailNotifier: Closing done channel to signal shutdown")
			close(done)
			log.Debug("emailNotifier: done channel closed, returning")
			return

		case responseDetails := <-incoming:

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

			t := time.NewTimer(sendTimeout)
			defer t.Stop()

			select {

			case <-t.C:

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

}

// StartNotifyMgr receives echoHandlerResponse values from a receive-only
// incoming queue of echoHandlerResponse values and sends notifications to any
// enabled service (e.g., Microsoft Teams).
func StartNotifyMgr(ctx context.Context, cfg *config.Config, notifyWorkQueue <-chan echoHandlerResponse, done chan<- struct{}) {

	log.Debug("StartNotifyMgr: Running")

	// Create separate, buffered channels to hand-off echoHandlerResponse
	// values for processing for each service, e.g., one channel for Microsoft
	// Teams outgoing notifications, another for email and so on. Buffered
	// channels are used both to enable async tasks and to provide a means of
	// monitoring the number of items queued for each channel; unbuffered
	// channels have a queue depth (and thus length) of 0.
	teamsNotifyWorkQueue := make(chan echoHandlerResponse, 10)
	teamsNotifyResultQueue := make(chan NotifyResult, 10)
	teamsNotifyDone := make(chan struct{})

	emailNotifyWorkQueue := make(chan echoHandlerResponse, 10)
	emailNotifyResultQueue := make(chan NotifyResult, 10)
	emailNotifyDone := make(chan struct{})

	if !cfg.NotifyTeams() && !cfg.NotifyEmail() {
		log.Debug("StartNotifyMgr: Teams and email notifications not requested, not starting notifier goroutines")
		// NOTE: Do not return/exit here.
		//
		// We cannot return/exit the function here because StartNotifyMgr HAS
		// to run in order to keep the notifyWorkQueue from filling up and
		// blocking other parts of this application that send messages to this
		// channel.
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
			teamsNotifyDone,
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
			emailNotifyDone,
		)
	}

	// Monitor queues and report stats for each
	if cfg.NotifyEmail() || cfg.NotifyTeams() {

		queuesToMonitor := []NotifyQueue{
			{
				Name:    "notifyWorkQueue",
				Channel: notifyWorkQueue,
			},
			{
				Name:    "emailNotifyWorkQueue",
				Channel: emailNotifyWorkQueue,
			},
			{
				Name:    "emailNotifyResultQueue",
				Channel: emailNotifyResultQueue,
			},
			{
				Name:    "teamsNotifyWorkQueue",
				Channel: teamsNotifyWorkQueue,
			},
			{
				Name:    "teamsNotifyResultQueue",
				Channel: teamsNotifyResultQueue,
			},
		}

		// print current queue items periodically
		go notifyQueueMonitor(ctx, config.NotifyQueueMonitorDelay, queuesToMonitor...)

	}

	for {

		select {

		// NOTE: This should ONLY ever be done when shutting down the entire
		// application, as otherwise goroutines associated with client
		// requests will likely hang, likely until client/server timeout
		// settings are reached
		case <-ctx.Done():
			// returning not to leak the goroutine
			ctxErr := ctx.Err()
			log.Debugf("StartNotifyMgr: Received Done signal: %v, shutting down ...", ctxErr.Error())

			evalResults := func(queueName string, result NotifyResult) {
				if result.Err != nil {
					log.Errorf("StartNotifyMgr: Error received from %s: %v", queueName, result.Err)
					return
				}
				log.Debugf("StartNotifyMgr: OK: non-error status received on %s: %v", queueName, result.Val)
			}

			// Process any waiting results before blocking and waiting
			// on final completion response from notifier goroutines
			if cfg.NotifyTeams() {
				log.Debug("Ranging over teamsNotifyResultQueue")
				for result := range teamsNotifyResultQueue {
					evalResults("teamsNotifyResultQueue", result)
				}

				log.Debug("StartNotifyMgr: Waiting on teamsNotifyDone")
				<-teamsNotifyDone
				log.Debug("StartNotifyMgr: Received from teamsNotifyDone")
			}

			if cfg.NotifyEmail() {
				log.Debug("Email notifications are enabled")
				log.Debug("Ranging over emailNotifyResultQueue")
				for result := range emailNotifyResultQueue {
					evalResults("emailNotifyResultQueue", result)
				}

				log.Debug("StartNotifyMgr: Waiting on emailNotifyDone")
				<-emailNotifyDone
				log.Debug("StartNotifyMgr: Received from emailNotifyDone")
			}

			log.Debug("StartNotifyMgr: Closing done channel")
			close(done)

			log.Debug("StartNotifyMgr: About to return")
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
					log.Debugf("StartNotifyMgr: Existing items in teamsNotifyWorkQueue: %d", len(teamsNotifyWorkQueue))
					log.Debug("StartNotifyMgr: Pending; placing responseDetails into teamsNotifyWorkQueue")
					teamsNotifyWorkQueue <- responseDetails
					log.Debug("StartNotifyMgr: Done; placed responseDetails into teamsNotifyWorkQueue")
					log.Debugf("StartNotifyMgr: Items now in teamsNotifyWorkQueue: %d", len(teamsNotifyWorkQueue))
				}()
			}

			if cfg.NotifyEmail() {
				log.Debug("StartNotifyMgr: Handing off responseDetails to emailNotifyWorkQueue")
				go func() {
					log.Debugf("StartNotifyMgr: Existing items in emailNotifyWorkQueue: %d", len(emailNotifyWorkQueue))
					log.Debug("StartNotifyMgr: Pending; placing responseDetails into emailNotifyWorkQueue")
					emailNotifyWorkQueue <- responseDetails
					log.Debug("StartNotifyMgr: Done; placed responseDetails into emailNotifyWorkQueue")
					log.Debugf("StartNotifyMgr: Items now in emailNotifyWorkQueue: %d", len(emailNotifyWorkQueue))
				}()
			}

			// default:
			// 	log.Debug("StartNotifyMgr: default case statement triggered")
		}

	}
}
