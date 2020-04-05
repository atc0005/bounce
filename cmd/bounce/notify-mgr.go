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

	// spin off goroutine to create and send Teams messages
	go func(ctx context.Context, incoming <-chan echoHandlerResponse, result chan<- error) {

		for {

			select {

			case <-ctx.Done():
				// returning not to leak the goroutine
				log.Debug("Received Done signal from context")
				return

			// Block waiting on input from notifyWorkQueue channel
			case responseDetails := <-incoming:

				log.Debugf("Request received by teams notification goroutine: %#v", responseDetails)

				// Send request details to Microsoft Teams if webhook URL set
				if cfg.NotifyTeams() {
					ourMessage := createMessage(responseDetails)
					if err := sendMessage(cfg.WebhookURL, ourMessage); err != nil {
						result <- fmt.Errorf("error occurred while trying to send message to Microsoft Teams: %w", err)
					}
				}

				// should a return be used here?

			}
		}

	}(ctx, teamsNotifyWorkQueue, teamsNotifyResultQueue)

	// spin off goroutine to create and send email messages
	go func(ctx context.Context, incoming <-chan echoHandlerResponse, result chan<- error) {

		for {

			select {

			case <-ctx.Done():
				// returning not to leak the goroutine
				log.Debug("Received Done signal from context")
				return

			// Block waiting on input from notifyWorkQueue channel
			case responseDetails := <-incoming:

				log.Debugf("Request received by email notification goroutine: %#v", responseDetails)

				// Send request details if enabled
				if !cfg.NotifyEmail() {
					errMsg := "Sending email is not currently enabled."
					log.Error(errMsg)
					result <- fmt.Errorf(errMsg)
				}

				// should a return be used here?
			}
		}

	}(ctx, emailNotifyWorkQueue, emailNotifyResultQueue)

	for {

		// Block waiting on input from notifyWorkQueue channel
		log.Debug("Waiting on input from notifyWorkQueue")
		responseDetails := <-notifyWorkQueue
		log.Debug("Input received from notifyWorkQueue")

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
