// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify"

	send2teams "github.com/atc0005/send2teams/teams"
)

func createMessage(clientRequest clientRequestDetails) goteamsnotify.MessageCard {

	// FIXME: This isn't an actual warning, just relying on color differences
	// during dev work for now.
	log.Debugf("createMessage: clientRequestDetails received: %#v", clientRequest)

	const ClientRequestErrorsRecorded = "Errors recorded for client request"
	const ClientRequestErrorsNotFound = "No errors recorded for client request"

	// FIXME: Pull this out as a separate helper function?
	// FIXME: Rework and offer upstream?
	addFactPair := func(msg *goteamsnotify.MessageCard, section *goteamsnotify.MessageCardSection, key string, values ...string) {

		if err := section.AddFactFromKeyValue(
			key,
			values...,
		); err != nil {

			// runtime.Caller(skip int) (pc uintptr, file string, line int, ok bool)
			_, file, line, ok := runtime.Caller(0)
			from := fmt.Sprintf("createMessage [file %s, line %d]:", file, line)
			if !ok {
				from = "createMessage:"
			}
			errMsg := fmt.Sprintf("%s error returned from attempt to add fact from key/value pair: %v", from, err)
			log.Errorf("%s %s", from, errMsg)
			msg.Text = msg.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
		}
	}

	// build MessageCard for submission
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = "Notification from " + config.MyAppName
	msgCard.Text = fmt.Sprintf(
		"%s request received on %s endpoint",
		send2teams.TryToFormatAsCodeSnippet(clientRequest.HTTPMethod),
		send2teams.TryToFormatAsCodeSnippet(clientRequest.EndpointPath),
	)

	/*
		Main Message Section
	*/

	// TODO: Is this needed?

	// mainMsgSection := goteamsnotify.NewMessageCardSection()
	// mainMsgSection.Title = "## Client request received"

	// if err := msgCard.AddSection(mainMsgSection); err != nil {
	// 	errMsg := fmt.Sprintf("\nError returned from attempt to add mainMsgSection: %v", err)
	// 	log.Error("createMessage: " + errMsg)
	// 	msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	// }

	// log.Info("This should show if the function is still running")

	/*
		Client Request Summary Section - General client request details
	*/

	clientRequestSummarySection := goteamsnotify.NewMessageCardSection()
	clientRequestSummarySection.Title = "## Client Request Summary"
	clientRequestSummarySection.StartGroup = true

	addFactPair(&msgCard, clientRequestSummarySection, "Received at", clientRequest.Datestamp)
	addFactPair(&msgCard, clientRequestSummarySection, "Endpoint path", clientRequest.EndpointPath)
	addFactPair(&msgCard, clientRequestSummarySection, "HTTP Method", clientRequest.HTTPMethod)
	addFactPair(&msgCard, clientRequestSummarySection, "Client IP Address", clientRequest.ClientIPAddress)

	if err := msgCard.AddSection(clientRequestSummarySection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientRequestSummarySection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Payload Section
	*/

	clientPayloadSection := goteamsnotify.NewMessageCardSection()
	clientPayloadSection.Title = "## Request body/payload"
	clientPayloadSection.StartGroup = true

	switch {
	case clientRequest.Body == "":
		log.Debugf("createMessage: Body is NOT defined, cannot use it to generate code block")
		clientPayloadSection.Text = send2teams.TryToFormatAsCodeSnippet("No request body was provided by client.")
	case clientRequest.Body != "":
		log.Debugf("createMessage: Body is defined, using it to generate code block")
		clientPayloadSection.Text = send2teams.TryToFormatAsCodeBlock(clientRequest.Body)
	}

	log.Debugf("createMessage: Body field contents: %v", clientRequest.Body)

	// FIXME: Remove this; only added for testing
	//clientPayloadSection.Text = ""

	if err := msgCard.AddSection(clientPayloadSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientPayloadSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Errors Section
	*/

	responseErrorsSection := goteamsnotify.NewMessageCardSection()
	responseErrorsSection.Title = "## Client Request errors"
	responseErrorsSection.StartGroup = true

	// Be optimistic to start with
	responseErrorsSection.Text = ClientRequestErrorsNotFound

	// Don't add this section if there are no errors to show
	if clientRequest.RequestError != "" {
		responseErrorsSection.Text = ""
		addFactPair(&msgCard, responseErrorsSection, "RequestError",
			send2teams.ConvertEOLToBreak(clientRequest.RequestError))
	}

	if clientRequest.BodyError != "" {
		responseErrorsSection.Text = ClientRequestErrorsRecorded
		addFactPair(&msgCard, responseErrorsSection, "BodyError",
			send2teams.ConvertEOLToBreak(clientRequest.BodyError))
	}

	if clientRequest.ContentTypeError != "" {
		responseErrorsSection.Text = ClientRequestErrorsRecorded
		addFactPair(&msgCard, responseErrorsSection, "ContentTypeError",
			send2teams.ConvertEOLToBreak(clientRequest.ContentTypeError))
	}

	if clientRequest.FormattedBodyError != "" {
		responseErrorsSection.Text = ClientRequestErrorsRecorded
		addFactPair(&msgCard, responseErrorsSection, "FormattedBodyError",
			send2teams.ConvertEOLToBreak(clientRequest.FormattedBodyError))
	}

	if err := msgCard.AddSection(responseErrorsSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add responseErrorsSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Headers Section
	*/

	clientRequestHeadersSection := goteamsnotify.NewMessageCardSection()
	clientRequestHeadersSection.StartGroup = true
	clientRequestHeadersSection.Title = "## Client Request Headers"

	clientRequestHeadersSection.Text = fmt.Sprintf(
		"%d client request headers provided",
		len(clientRequest.Headers),
	)

	// process client request headers

	for header, values := range clientRequest.Headers {
		for index, value := range values {
			// update value with code snippet formatting, assign back using
			// the available index value
			values[index] = send2teams.TryToFormatAsCodeSnippet(value)
		}
		addFactPair(&msgCard, clientRequestHeadersSection, header, values...)
	}

	if err := msgCard.AddSection(clientRequestHeadersSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientRequestHeadersSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Message Card Branding/Trailer Section
	*/

	trailerSection := goteamsnotify.NewMessageCardSection()
	trailerSection.StartGroup = true
	trailerSection.Text = send2teams.ConvertEOLToBreak(config.MessageTrailer())
	if err := msgCard.AddSection(trailerSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add trailerSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + send2teams.TryToFormatAsCodeSnippet(errMsg)
	}

	return msgCard
}

// define function/wrapper for sending details to Microsoft Teams
func sendMessage(
	ctx context.Context,
	webhookURL string,
	msgCard goteamsnotify.MessageCard,
	schedule time.Time,
	retries int,
	retriesDelay int,
) NotifyResult {

	// Note: We already do validation elsewhere, and the library call does
	// even more validation, but we can handle this obvious empty argument
	// problem directly
	if webhookURL == "" {
		return NotifyResult{
			Err: fmt.Errorf("sendMessage: webhookURL not defined, skipping message submission to Microsoft Teams channel"),
		}
	}

	// Take received schedule and use it to determine how long our timer
	// should be before we make our first attempt at sending a message to
	// Microsoft Teams
	nextSchedule := time.Until(schedule)

	notificationDelayTimer := time.NewTimer(nextSchedule)
	defer notificationDelayTimer.Stop()
	log.Debugf("sendMessage: notificationDelayTimer created at %v with duration %v",
		time.Now().Format("15:04:05"),
		config.NotifyMgrTeamsNotificationDelay,
	)

	log.Debug("sendMessage: Waiting for either context or notificationDelayTimer to expire before sending notification")

	select {
	case <-ctx.Done():
		ctxErr := ctx.Err()
		msg := NotifyResult{
			Val: fmt.Sprintf("sendMessage: Received Done signal at %v: %v, shutting down",
				time.Now().Format("15:04:05"),
				ctxErr.Error(),
			),
		}
		log.Debug(msg.Val)
		return msg

	// Delay between message submission attempts; this will *always*
	// delay, regardless of whether the attempt is the first one or not
	case <-notificationDelayTimer.C:

		log.Debugf("sendMessage: Waited %v before notification attempt at %v",
			config.NotifyMgrTeamsNotificationDelay,
			time.Now().Format("15:04:05"),
		)

		ctxExpires, ctxExpired := ctx.Deadline()
		if ctxExpired {
			log.Debugf("sendMessage: WaitTimeout context expires at: %v", ctxExpires.Format("15:04:05"))
		}

		// check to see if context has expired during our delay
		if ctx.Err() != nil {
			msg := NotifyResult{
				Val: fmt.Sprintf(
					"sendMessage: context expired or cancelled: %v, attempting to abort message submission",
					ctx.Err().Error(),
				),
			}

			log.Debug(msg.Val)

			return msg
		}

		// Submit message card, retry submission if needed up to specified number
		// of retry attempts.
		if err := send2teams.SendMessage(ctx, webhookURL, msgCard, retries, retriesDelay); err != nil {
			errMsg := NotifyResult{
				Err: fmt.Errorf("sendMessage: ERROR: Failed to submit message to Microsoft Teams: %v", err),
			}
			log.Error(errMsg.Err.Error())
			return errMsg
		}

		successMsg := NotifyResult{
			Val: "sendMessage: Message successfully sent to Microsoft Teams",
		}

		// Note success for potential troubleshooting
		log.Debug(successMsg.Val)

		return successMsg

	}

}
