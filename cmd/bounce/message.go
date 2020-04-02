// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/atc0005/bounce/config"

	// use our fork for now until recent work can be submitted for inclusion
	// in the upstream project
	goteamsnotify "github.com/atc0005/go-teams-notify"

	send2teams "github.com/atc0005/send2teams/teams"
)

func createMessage(responseDetails echoHandlerResponse) goteamsnotify.MessageCard {

	// FIXME: This isn't an actual warning, just relying on color differences
	// during dev work for now.
	log.Warnf("echoHandlerResponse received: %#v", responseDetails)

	// build MessageCard for submission
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = "Notification from " + config.MyAppName
	msgCard.Text = fmt.Sprintf(
		"%s request received on %s endpoint",
		goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.HTTPMethod),
		goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.EndpointPath),
	)

	/*
		Main Message Section
	*/

	// TODO: Is this needed?

	// mainMsgSection := goteamsnotify.NewMessageCardSection()
	// mainMsgSection.Title = "## Client request received"

	// if err := msgCard.AddSection(mainMsgSection); err != nil {
	// 	errMsg := fmt.Sprintf("\nError returned from attempt to add mainMsgSection: %v", err)
	// 	log.Error(errMsg)
	// 	msgCard.Text = msgCard.Text + "\n\n" + goteamsnotify.TryToFormatAsCodeSnippet(errMsg)
	// }

	// log.Info("This should show if the function is still running")

	/*
		Client Request Payload Section
	*/

	clientPayloadSection := goteamsnotify.NewMessageCardSection()
	clientPayloadSection.Title = "## Request body/payload"
	clientPayloadSection.StartGroup = true

	switch {
	case responseDetails.Body == "":
		log.Debugf("Body is NOT defined, cannot use it to generate code block")
		clientPayloadSection.Text = goteamsnotify.TryToFormatAsCodeSnippet("No request body was provided by client.")
	case responseDetails.Body != "":
		log.Debugf("Body is defined, using it to generate code block")
		clientPayloadSection.Text = goteamsnotify.TryToFormatAsCodeBlock(responseDetails.Body)
	}

	log.Debugf("Body field contents: %v", responseDetails.Body)

	// FIXME: Remove this; only added for testing
	//clientPayloadSection.Text = ""

	if err := msgCard.AddSection(clientPayloadSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientPayloadSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + goteamsnotify.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Summary Section - General client request details
	*/

	clientRequestSummarySection := goteamsnotify.NewMessageCardSection()
	clientRequestSummarySection.Title = "## Client Request Summary"
	clientRequestSummarySection.StartGroup = true

	clientRequestSummarySection.AddFactFromKeyValue(
		"Received at",
		goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.Datestamp),
	)

	clientRequestSummarySection.AddFactFromKeyValue(
		"Endpoint path",
		goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.EndpointPath),
	)

	clientRequestSummarySection.AddFactFromKeyValue(
		"HTTP Method",
		goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.HTTPMethod),
	)

	clientRequestSummarySection.AddFactFromKeyValue(
		"Client IP Address",
		goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.ClientIPAddress),
	)

	if err := msgCard.AddSection(clientRequestSummarySection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientRequestSummarySection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + goteamsnotify.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Errors Section
	*/

	responseErrorsSection := goteamsnotify.NewMessageCardSection()
	responseErrorsSection.Title = "## Client Request errors"
	responseErrorsSection.StartGroup = true

	// Be optimistic to start with
	responseErrorsSection.Text = "No errors recorded for client request."

	// Don't add this section if there are no errors to show
	if responseDetails.RequestError != "" {

		responseErrorsSection.Text = ""
		responseErrorsSection.AddFactFromKeyValue(
			"RequestError",
			//goteamsnotify.TryToFormatAsCodeSnippet(responseDetails.RequestError),
			send2teams.ConvertEOLToBreak(responseDetails.RequestError),
		)
	}

	if responseDetails.BodyError != "" {

		responseErrorsSection.Text = "Errors recorded for client request"
		responseErrorsSection.AddFactFromKeyValue(
			"BodyError",
			send2teams.ConvertEOLToBreak(responseDetails.BodyError),
		)
	}

	if responseDetails.ContentTypeError != "" {

		responseErrorsSection.Text = "Errors recorded for client request"
		responseErrorsSection.AddFactFromKeyValue(
			"ContentTypeError",
			send2teams.ConvertEOLToBreak(responseDetails.ContentTypeError),
		)
	}

	if responseDetails.FormattedBodyError != "" {

		responseErrorsSection.Text = "Errors recorded for client request"
		responseErrorsSection.AddFactFromKeyValue(
			"FormattedBodyError",
			send2teams.ConvertEOLToBreak(responseDetails.FormattedBodyError),
		)

	}

	if err := msgCard.AddSection(responseErrorsSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add responseErrorsSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + goteamsnotify.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Client Request Headers Section
	*/

	clientRequestHeadersSection := goteamsnotify.NewMessageCardSection()
	clientRequestHeadersSection.StartGroup = true
	clientRequestHeadersSection.Title = "## Client Request Headers"

	clientRequestHeadersSection.Text = fmt.Sprintf(
		"%d client request headers provided",
		len(responseDetails.Headers),
	)

	// process client request headers

	for header, values := range responseDetails.Headers {
		for index, value := range values {
			// update value with code snippet formatting, assign back using
			// the available index value
			values[index] = goteamsnotify.TryToFormatAsCodeSnippet(value)
		}
		clientRequestHeadersSection.AddFactFromKeyValue(header, values...)
	}

	if err := msgCard.AddSection(clientRequestHeadersSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientRequestHeadersSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + goteamsnotify.TryToFormatAsCodeSnippet(errMsg)
	}

	/*
		Message Card Branding/Trailer Section
	*/

	trailerSection := goteamsnotify.NewMessageCardSection()
	trailerSection.StartGroup = true
	trailerSection.Text = send2teams.ConvertEOLToBreak(config.MessageTrailer())
	if err := msgCard.AddSection(trailerSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add trailerSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + goteamsnotify.TryToFormatAsCodeSnippet(errMsg)
	}

	return msgCard
}

// define function/wrapper for sending details to Microsoft Teams
func sendMessage(webhookURL string, msgCard goteamsnotify.MessageCard) error {

	if webhookURL == "" {
		log.Debug("webhookURL not defined, skipping message submission to Microsoft Teams channel")
	}

	// Submit message card
	if err := send2teams.SendMessage(webhookURL, msgCard); err != nil {
		errMsg := fmt.Errorf("ERROR: Failed to submit message to Microsoft Teams: %v", err)
		log.Error(errMsg.Error())
		return errMsg
	}

	// NOTE: Unscientific testing showed a MASSIVE difference in
	// response times when launching this in a goroutine. We'll
	// need to find a way to communicate *back* to the caller
	// the results of the goroutine, otherwise we are not
	// able to properly handle errors.

	// go func() error {

	// 	// Submit message card
	// 	if err := send2teams.SendMessage(webhookURL, msgCard); err != nil {
	// 		errMsg := fmt.Errorf("ERROR: Failed to submit message to Microsoft Teams: %v", err)
	// 		log.Error(errMsg.Error())
	// 		return errMsg
	// 	}

	// 	// Emit basic success message
	// 	log.Info("Message successfully sent to Microsoft Teams")
	// 	return nil

	// }()

	// Emit basic success message
	log.Info("Message successfully sent to Microsoft Teams")

	return nil

}
