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

	// build MessageCard for submission
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = "Client request received"
	msgCard.Text = "Our first test from bounce!"

	mainMsgSection := goteamsnotify.NewMessageCardSection()
	// mainMsgSection.Title = "Client request received (mainMsgSection.Title)"

	if err := msgCard.AddSection(mainMsgSection); err != nil {
		errMsg := fmt.Sprintf("\nError returned from attempt to add mainMsgSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n" + errMsg
	}

	log.Info("This should show if the function is still running")

	clientPayloadSection := goteamsnotify.NewMessageCardSection()
	clientPayloadSection.Title = "Client request details"

	// JSON payload if available
	// echoHandlerResponse struct {
	// 	Datestamp          string
	// 	EndpointPath       string
	// 	HTTPMethod         string
	// 	ClientIPAddress    string
	// 	Headers            http.Header
	// 	Body               string
	// 	BodyError          string
	// 	FormattedBody      string
	// 	FormattedBodyError string
	// 	RequestError       string
	// 	ContentTypeError   string
	// }

	// 	Request received: 2020-04-02T07:30:11-05:00
	// Endpoint path requested by client: /api/v1/echo
	// HTTP Method used by client: GET
	// Client IP Address: 127.0.0.1:61452

	// Headers:

	//   * Accept: */*
	//   * User-Agent: curl/7.68.0

	switch {
	case responseDetails.Body == "":
		log.Debugf("Body is NOT defined, cannot use it to generate code block")
		clientPayloadSection.Text = "No request body was provided by client."
	case responseDetails.Body != "":
		log.Debugf("Body is defined, using it to generate code block")
		codeBlock, err := goteamsnotify.FormatAsCodeBlock(responseDetails.Body)
		if err != nil {
			// Should be something like this:
			// "No request body was provided by client."
			clientPayloadSection.Text = err.Error()
		}
		clientPayloadSection.Text = codeBlock
	}

	log.Debugf("Body field contents: %v", responseDetails.Body)

	// FIXME: Remove this; only added for testing
	//clientPayloadSection.Text = ""

	if err := msgCard.AddSection(clientPayloadSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add clientPayloadSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + errMsg
	}

	trailerSection := goteamsnotify.NewMessageCardSection()
	trailerSection.Text = send2teams.ConvertEOLToBreak(config.Branding())
	if err := msgCard.AddSection(trailerSection); err != nil {
		errMsg := fmt.Sprintf("Error returned from attempt to add trailerSection: %v", err)
		log.Error(errMsg)
		msgCard.Text = msgCard.Text + "\n\n" + errMsg
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

	// Emit basic success message
	log.Info("Message successfully sent to Microsoft Teams")

	return nil

}
