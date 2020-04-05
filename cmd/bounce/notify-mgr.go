package main

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"
)

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

	// Block waiting on input from notifyWorkQueue channel
	responseDetails := <-notifyWorkQueue

	// spin off goroutine to create and send Teams messages
	go func(incoming <-chan echoHandlerResponse, result chan<- error) {

		// TODO: setup infinite loop to process incoming items

		// Send request details to Microsoft Teams if webhook URL set
		if cfg.NotifyTeams() {
			ourMessage := createMessage(responseDetails)
			if err := sendMessage(cfg.WebhookURL, ourMessage); err != nil {
				result <- fmt.Errorf("error occurred while trying to send message to Microsoft Teams: %w", err)
			}
		}

		result <- nil
	}(teamsNotifyWorkQueue, teamsNotifyResultQueue)

	// spin off goroutine to create and send email messages
	go func(incoming <-chan echoHandlerResponse, result chan<- error) {

		// TODO: setup infinite loop to process incoming items

		// Send request details if enabled
		if !cfg.NotifyEmail() {
			errMsg := "Sending email is not currently enabled."
			log.Error(errMsg)
			result <- fmt.Errorf(errMsg)
		}

		// this shouldn't be reached
		result <- nil
	}(emailNotifyWorkQueue, emailNotifyResultQueue)

	// FIXME: Is this for loop dedicated to just receiving values? If so, we
	// should not insert any statements that sent values down a channel ...
	for {
		select {
		case <-ctx.Done():
			// returning not to leak the goroutine
			log.Debug("Received Done signal from context")
			return

		// Attempt to send response details to goroutine responsible for
		// generating Microsoft Teams messages
		case teamsNotifyWorkQueue <- responseDetails:
			log.Debug("Handed off responseDetails to teamsNotifyWorkQueue")
			// success; now what?

		case err := <-teamsNotifyResultQueue:
			// do something based on success or failure sending to Teams
			if err != nil {
				log.Error(err.Error())
			}

		// Attempt to send response details to goroutine responsible for
		// generating email messages
		case emailNotifyWorkQueue <- responseDetails:
			log.Debug("Handed off responseDetails to emailNotifyWorkQueue")
			// success; now what?

		case err := <-emailNotifyResultQueue:
			// do something based on success or failure sending to Teams
			if err != nil {
				log.Error(err.Error())
			}
		}

	}

}

// At least one goroutine running an infinite for loop with a cancel context.
// Process queue, wait on queue items.

// Would channels function as queues?

// Q: How to make the goroutine pause while waiting on items to be added to
// the queue?

// A: select statement

// ---

// main calls a function that launches an infinite loop goroutine with a
// cancel context. Wait on incoming messages, spin off separate goroutine when
// needed to handle sending the message. This "notification manager" goroutine
// can handle errors from message send operations.

// Potentially it could also handle email notifications too?

// Incoming email channel and incoming teams channel? Perhaps this manager
// goroutine can use a single incoming channel to receive the responseDetails
// (perhaps ResponseDetails) type and also check flag settings. If Teams
// support enabled, then send message via Teams. If email support enabled,
// send email. Both send msg types would be by new goroutines.
