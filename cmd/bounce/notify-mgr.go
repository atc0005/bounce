package main

import "context"

// StartNotifyMgr receives echoHandlerResponse values from a receive-only
// incoming queue of echoHandlerResponse values and sends notifications to any
// enabled service (e.g., Microsoft Teams).
func StartNotifyMgr(ctx context.Context, notifyWorkQueue <-chan echoHandlerResponse) {

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
	//emailNotifyWorkQueue := make(chan echoHandlerResponse)

	// Block waiting on input from notifyWorkQueue channel
	responseDetails := <-notifyWorkQueue

	// FIXME: Is this for loop dedicated to just receiving values? If so, we
	// should not insert any statements that sent values down a channel ...
	for {
		select {
		case <-ctx.Done():
			// returning not to leak the goroutine
			return

		// Attempt to send response details to goroutine responsible for
		// generating Microsoft Teams messages
		case teamsNotifyWorkQueue <- responseDetails:
			// success; now what?
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
