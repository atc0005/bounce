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

	// Val is the non-error condition message to return from a notification
	// operation
	Val string

	// Err is the error condition message to return from a notification
	// operation
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

	for {

		t := time.NewTimer(delay)

		// log.Debug("notifyQueueMonitor: Starting loop")

		// block until:
		//	- context cancellation
		//	- timer fires
		select {
		case <-ctx.Done():
			t.Stop()
			log.Debugf(
				"notifyQueueMonitor: Received Done signal: %v, shutting down ...",
				ctx.Err().Error(),
			)
			return

		case <-t.C:

			// log.Debug("notifyQueueMonitor: Timer fired")

			// NOTE: Not needed since the channel is already drained as a
			// result of the case statement triggering and draining the
			// channel
			// t.Stop()

			var itemsFound bool
			//log.Debugf("Length of queues: %d", len(queues))
			for _, notifyQueue := range notifyQueues {

				var queueLength int
				switch queue := notifyQueue.Channel.(type) {

				// FIXME: Is there a generic way to match any channel type
				// here in order to calculate the length?
				case chan clientRequestDetails:
					queueLength = len(queue)

				case <-chan clientRequestDetails:
					queueLength = len(queue)

				case chan NotifyResult:
					queueLength = len(queue)

				default:
					log.Warn("Default case triggered (this should not happen")
					log.Warnf("Name of channel: %s", notifyQueue.Name)

				}

				// Show stats only for queues with content
				if queueLength > 0 {
					itemsFound = true
					log.Debugf("notifyQueueMonitor: %d items in %s",
						queueLength, notifyQueue.Name)
					log.Debugf("notifyQueueMonitor: %d goroutines running",
						runtime.NumGoroutine())
					continue
				}

			}

			if !itemsFound {
				log.Debugf("notifyQueueMonitor: 0 items queued, %d goroutines running",
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
	incoming <-chan clientRequestDetails,
	notifyMgrResultQueue chan<- NotifyResult,
	done chan<- struct{},
) {

	// TODO: Replace config package constant references with function parameters?

	log.Debug("teamsNotifier: Running")

	// used by goroutines called by this function to return results
	ourResultQueue := make(chan NotifyResult)

	// We need to account for multiple factors when we set a complete
	// timeout for sending messages to Teams:
	//
	// - the base timeout value for a single message submission attempt
	// - the delay we are enforcing between message submission attempts
	// - the total number of retries allowed
	// - the delay between retry attempts
	timeoutValue := (config.NotifyMgrTeamsTimeout +
		config.NotifyMgrTeamsNotificationDelay +
		time.Duration(retriesDelay)) * time.Duration(retries)

	for {

		select {

		case <-ctx.Done():

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

		case clientRequest := <-incoming:

			// FIXME: The timer handling needs additional testing (very little has been done so far)
			// one-time events, have to recreate timer on each iteration
			timeoutTimer := time.NewTimer(timeoutValue)
			log.Debugf("teamsNotifier: timeoutTimer created with duration %v", timeoutValue)

			// TODO: Do we need to also check context state here?
			//
			// i.e.g, if there is a message waiting *and* ctx.Done() case
			// statements are both valid, either path could be taken. If this
			// one is taken, then the message send timeout will be the only
			// thing forcing the attempt to loop back around and trigger the
			// ctx.Done() path, but only if this one isn't taken again by the
			// random case selection logic

			log.Debugf("teamsNotifier: Request received at %v: %#v",
				time.Now(), clientRequest)

			log.Debug("teamsNotifier: Checking context to determine whether we should proceed")
			if ctx.Err() != nil {
				log.Debug("teamsNotifier: context has been cancelled, aborting notification attempt")

				// stop all timers
				timeoutTimer.Stop()
				continue
			}
			log.Debug("teamsNotifier: context not cancelled, proceeding with notification attempt")

			// launch task in separate goroutine
			log.Debug("teamsNotifier: Launching message creation/submission in separate goroutine")
			go func(ctx context.Context, webhookURL string, clientRequest clientRequestDetails, resultQueue chan<- NotifyResult) {
				ourMessage := createMessage(clientRequest)
				resultQueue <- sendMessage(ctx, webhookURL, ourMessage, retries, retriesDelay)
				return
			}(ctx, webhookURL, clientRequest, ourResultQueue)

			select {

			// timeout for the entire message submission
			// if this occurs we just move on to the next message
			case <-timeoutTimer.C:

				result := NotifyResult{
					Err: fmt.Errorf(
						"teamsNotifier: Timeout reached at %v (%v) after %d attempt to send Microsoft Teams notification",
						time.Now(),
						sendTimeout,
						retries+1,
					),
				}
				log.Debug(result.Err.Error())
				notifyMgrResultQueue <- result
				continue

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
func emailNotifier(ctx context.Context, sendTimeout time.Duration, incoming <-chan clientRequestDetails, notifyMgrResultQueue chan<- NotifyResult, done chan<- struct{}) {

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

		case clientRequest := <-incoming:

			log.Debugf("emailNotifier: Request received: %#v", clientRequest)

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

// StartNotifyMgr receives clientRequestDetails values from a receive-only
// incoming queue of clientRequestDetails values and sends notifications to any
// enabled service (e.g., Microsoft Teams).
// FIXME: Tweak the description for this function; it seems to have some stutter
func StartNotifyMgr(ctx context.Context, cfg *config.Config, notifyWorkQueue <-chan clientRequestDetails, done chan<- struct{}) {

	log.Debug("StartNotifyMgr: Running")

	// Create separate, buffered channels to hand-off clientRequestDetails
	// values for processing for each service, e.g., one channel for Microsoft
	// Teams outgoing notifications, another for email and so on. Buffered
	// channels are used both to enable async tasks and to provide a means of
	// monitoring the number of items queued for each channel; unbuffered
	// channels have a queue depth (and thus length) of 0.
	teamsNotifyWorkQueue := make(chan clientRequestDetails, 5)
	teamsNotifyResultQueue := make(chan NotifyResult, 5)
	teamsNotifyDone := make(chan struct{})

	emailNotifyWorkQueue := make(chan clientRequestDetails, 5)
	emailNotifyResultQueue := make(chan NotifyResult, 5)
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
		log.Debug("StartNotifyMgr: Teams notifications enabled")
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
		log.Debug("StartNotifyMgr: Email notifications enabled")
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
				log.Debug("StartNotifyMgr: Teams notifications are enabled, shutting down teamsNotifier")

				log.Debug("StartNotifyMgr: Ranging over teamsNotifyResultQueue")
				for result := range teamsNotifyResultQueue {
					evalResults("teamsNotifyResultQueue", result)
				}

				log.Debug("StartNotifyMgr: Waiting on teamsNotifyDone")
				select {
				case <-teamsNotifyDone:
					log.Debug("StartNotifyMgr: Received from teamsNotifyDone")
				case <-time.After(config.NotifyMgrServicesShutdownTimeout):
					log.Debug("StartNotifyMgr: Timeout occurred while waiting for teamsNotifyDone")
					log.Debug("StartNotifyMgr: Proceeding with shutdown")
				}

			}

			if cfg.NotifyEmail() {
				log.Debug("StartNotifyMgr: Email notifications are enabled, shutting down emailNotifier")

				log.Debug("StartNotifyMgr: Ranging over emailNotifyResultQueue")
				for result := range emailNotifyResultQueue {
					evalResults("emailNotifyResultQueue", result)
				}

				log.Debug("StartNotifyMgr: Waiting on emailNotifyDone")
				select {
				case <-emailNotifyDone:
					log.Debug("StartNotifyMgr: Received from emailNotifyDone")
				case <-time.After(config.NotifyMgrServicesShutdownTimeout):
					log.Debug("StartNotifyMgr: Timeout occurred while waiting for emailNotifyDone")
					log.Debug("StartNotifyMgr: Proceeding with shutdown")
				}

			}

			log.Debug("StartNotifyMgr: Closing done channel")
			close(done)

			log.Debug("StartNotifyMgr: About to return")
			return

		case clientRequest := <-notifyWorkQueue:

			log.Debug("StartNotifyMgr: Input received from notifyWorkQueue")

			// If we don't have *any* notifications enabled we will just
			// discard the item we have pulled from the channel
			if !cfg.NotifyEmail() && !cfg.NotifyTeams() {
				log.Debug("StartNotifyMgr: Notifications are not currently enabled; ignoring notification request")
				continue
			}

			if cfg.NotifyTeams() {
				log.Debug("StartNotifyMgr: Creating new goroutine to place clientRequest into teamsNotifyWorkQueue")
				go func() {
					log.Debugf("StartNotifyMgr: Existing items in teamsNotifyWorkQueue: %d", len(teamsNotifyWorkQueue))
					log.Debug("StartNotifyMgr: Pending; placing clientRequest into teamsNotifyWorkQueue")
					teamsNotifyWorkQueue <- clientRequest
					log.Debug("StartNotifyMgr: Done; placed clientRequest into teamsNotifyWorkQueue")
					log.Debugf("StartNotifyMgr: Items now in teamsNotifyWorkQueue: %d", len(teamsNotifyWorkQueue))
				}()
			}

			if cfg.NotifyEmail() {
				log.Debug("StartNotifyMgr: Creating new goroutine to place clientRequest in emailNotifyWorkQueue")
				go func() {
					log.Debugf("StartNotifyMgr: Existing items in emailNotifyWorkQueue: %d", len(emailNotifyWorkQueue))
					log.Debug("StartNotifyMgr: Pending; placing clientRequest into emailNotifyWorkQueue")
					emailNotifyWorkQueue <- clientRequest
					log.Debug("StartNotifyMgr: Done; placed clientRequest into emailNotifyWorkQueue")
					log.Debugf("StartNotifyMgr: Items now in emailNotifyWorkQueue: %d", len(emailNotifyWorkQueue))
				}()
			}

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

		}

	}
}
