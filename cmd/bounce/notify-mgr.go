package main

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
