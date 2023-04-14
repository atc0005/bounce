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
	"github.com/atc0005/bounce/internal/config"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

func createMessage(clientRequest clientRequestDetails) *messagecard.MessageCard {

	log.Debugf("createMessage: clientRequestDetails received: %#v", clientRequest)

	const ClientRequestErrorsRecorded = "Errors recorded for client request"
	const ClientRequestErrorsNotFound = "No errors recorded for client request"

	// FIXME: Pull this out as a separate helper function?
	// FIXME: Rework and offer upstream?
	addFactPair := func(msg *messagecard.MessageCard, section *messagecard.Section, key string, values ...string) {

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
			msg.Text = msg.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
		}
	}

	// build MessageCard for submission
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = "Notification from " + config.MyAppName
	msgCard.Text = fmt.Sprintf(
		"%s request received on %s endpoint",
		messagecard.TryToFormatAsCodeSnippet(clientRequest.HTTPMethod),
		messagecard.TryToFormatAsCodeSnippet(clientRequest.EndpointPath),
	)

	/*
		Client Request Summary Section - General client request details
	*/

	clientRequestSummarySection := messagecard.NewSection()
	clientRequestSummarySection.Title = "## Client Request Summary"
	clientRequestSummarySection.StartGroup = true

	addFactPair(msgCard, clientRequestSummarySection, "Received at", clientRequest.Datestamp)
	addFactPair(msgCard, clientRequestSummarySection, "Endpoint path", clientRequest.EndpointPath)
	addFactPair(msgCard, clientRequestSummarySection, "HTTP Method", clientRequest.HTTPMethod)
	addFactPair(msgCard, clientRequestSummarySection, "Client IP Address", clientRequest.ClientIPAddress)

	if err := msgCard.AddSection(clientRequestSummarySection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientRequestSummarySection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Payload Section
	*/

	clientPayloadSection := messagecard.NewSection()
	clientPayloadSection.Title = "## Request body/payload"
	clientPayloadSection.StartGroup = true

	switch {
	case clientRequest.Body == "":
		log.Debugf("createMessage: Body is NOT defined, cannot use it to generate code block")
		clientPayloadSection.Text = messagecard.TryToFormatAsCodeSnippet("No request body was provided by client.")
	case clientRequest.Body != "":
		log.Debugf("createMessage: Body is defined, using it to generate code block")
		clientPayloadSection.Text = messagecard.TryToFormatAsCodeBlock(clientRequest.Body)
	}

	log.Debugf("createMessage: Body field contents: %v", clientRequest.Body)

	if err := msgCard.AddSection(clientPayloadSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientPayloadSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Errors Section
	*/

	responseErrorsSection := messagecard.NewSection()
	responseErrorsSection.Title = "## Client Request errors"
	responseErrorsSection.StartGroup = true

	// Be optimistic to start with
	responseErrorsSection.Text = ClientRequestErrorsNotFound

	if clientRequest.RequestError != "" {
		responseErrorsSection.Text = ""
		addFactPair(msgCard, responseErrorsSection, "RequestError",
			messagecard.ConvertEOLToBreak(clientRequest.RequestError))
	}

	if clientRequest.BodyError != "" {
		responseErrorsSection.Text = ClientRequestErrorsRecorded
		addFactPair(msgCard, responseErrorsSection, "BodyError",
			messagecard.ConvertEOLToBreak(clientRequest.BodyError))
	}

	if clientRequest.ContentTypeError != "" {
		responseErrorsSection.Text = ClientRequestErrorsRecorded
		addFactPair(msgCard, responseErrorsSection, "ContentTypeError",
			messagecard.ConvertEOLToBreak(clientRequest.ContentTypeError))
	}

	if clientRequest.FormattedBodyError != "" {
		responseErrorsSection.Text = ClientRequestErrorsRecorded
		addFactPair(msgCard, responseErrorsSection, "FormattedBodyError",
			messagecard.ConvertEOLToBreak(clientRequest.FormattedBodyError))
	}

	if err := msgCard.AddSection(responseErrorsSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add responseErrorsSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Headers Section
	*/

	clientRequestHeadersSection := messagecard.NewSection()
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
			values[index] = messagecard.TryToFormatAsCodeSnippet(value)
		}
		addFactPair(msgCard, clientRequestHeadersSection, header, values...)
	}

	if err := msgCard.AddSection(clientRequestHeadersSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientRequestHeadersSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Message Card Branding/Trailer Section
	*/

	trailerSection := messagecard.NewSection()
	trailerSection.StartGroup = true
	trailerSection.Text = messagecard.ConvertEOLToBreak(config.MessageTrailer())
	if err := msgCard.AddSection(trailerSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add trailerSection: %v", err)
		log.Error("createMessage: " + errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + messagecard.TryToFormatAsCodeSnippet(errMsg)
	}

	return msgCard
}

// define function/wrapper for sending details to Microsoft Teams
func sendMessage(
	ctx context.Context,
	webhookURL string,
	msgCard *messagecard.MessageCard,
	schedule time.Time,
	retries int,
	retriesDelay int,
) NotifyResult {

	// Note: We already do validation elsewhere, and the library call does
	// even more validation, but we can handle this obvious empty argument
	// problem directly
	if webhookURL == "" {
		return NotifyResult{
			Err:     fmt.Errorf("sendMessage: webhookURL not defined, skipping message submission to Microsoft Teams channel"),
			Success: false,
		}
	}

	log.Debugf("sendMessage: Time now is %v", time.Now().Format("15:04:05"))
	log.Debugf("sendMessage: Notification scheduled for: %v", schedule.Format("15:04:05"))

	// Set delay timer to meet received notification schedule. This helps
	// ensure that we delay the appropriate amount of time before we make our
	// first attempt at sending a message to Microsoft Teams.
	notificationDelay := time.Until(schedule)

	notificationDelayTimer := time.NewTimer(notificationDelay)
	defer notificationDelayTimer.Stop()
	log.Debugf("sendMessage: notificationDelayTimer created at %v with duration %v",
		time.Now().Format("15:04:05"),
		notificationDelay,
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
			Success: false,
		}
		log.Debug(msg.Val)
		return msg

	// Delay between message submission attempts; this will *always*
	// delay, regardless of whether the attempt is the first one or not
	case <-notificationDelayTimer.C:

		log.Debugf("sendMessage: Waited %v before notification attempt at %v",
			notificationDelay,
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
					"sendMessage: context expired or cancelled at %v: %v, attempting to abort message submission",
					time.Now().Format("15:04:05"),
					ctx.Err().Error(),
				),
				Success: false,
			}

			log.Debug(msg.Val)

			return msg
		}

		// Create Microsoft Teams client
		mstClient := goteamsnotify.NewTeamsClient()

		// Submit message card using Microsoft Teams client, retry submission
		// if needed up to specified number of retry attempts.
		if err := mstClient.SendWithRetry(ctx, webhookURL, msgCard, retries, retriesDelay); err != nil {
			errMsg := NotifyResult{
				Err: fmt.Errorf(
					"sendMessage: ERROR: Failed to submit message to Microsoft Teams at %v: %w",
					time.Now().Format("15:04:05"),
					err,
				),
				Success: false,
			}
			log.Error(errMsg.Err.Error())
			return errMsg
		}

		successMsg := NotifyResult{
			Val: fmt.Sprintf(
				"sendMessage: Message successfully sent to Microsoft Teams at %v",
				time.Now().Format("15:04:05"),
			),
			Success: true,
		}

		// Note success for potential troubleshooting
		log.Debug(successMsg.Val)

		return successMsg

	}

}
