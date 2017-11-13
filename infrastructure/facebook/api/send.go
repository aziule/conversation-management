package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	ErrCouldNotMarshalJson = errors.New("Could not marshal JSON object")
)

// SendTextToUser is the FacebookApi's interface method responsible for sending a 1-to-1 message to a user
func (api *FacebookApi) SendTextToUser(recipientId, text string) error {
	url := api.getSendTextUrl()

	object := newTextToUserEnvelope(recipientId, text)
	jsonObject, err := json.Marshal(object)

	if err != nil {
		log.WithFields(log.Fields{
			"recipientId": recipientId,
			"text":        text,
		}).Infof("Could not send the message to the user due to a JSON marshal issue: %s", err)
		return ErrCouldNotMarshalJson
	}

	request, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonObject))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.WithFields(log.Fields{
			"url":        url.String(),
			"jsonObject": string(jsonObject),
		}).Infof("Could not create a new request: %s", err)
		return err
	}

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		log.Infof("Failed to send the request: %s", err)
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Infof("Failed to read the response's body: %s", err)
	}

	// @todo: handle errors
	fmt.Println(string(body))

	return nil
}

// recipientEnvelope is the envelope for a recipient
type recipientEnvelope struct {
	Id string `json:"id"`
}

// messageEnvelope represents the envelope for a message with text
type messageEnvelope struct {
	Text string `json:"text"`
}

// textToUserEnvelope is the JSON envelope that needs to be sent
type textToUserEnvelope struct {
	Recipient *recipientEnvelope `json:"recipient"`
	Message   *messageEnvelope   `json:"message"`
}

// newRecipientEnvelope is the constructor for a recipientEnvelope
func newRecipientEnvelope(recipientId string) *recipientEnvelope {
	return &recipientEnvelope{
		Id: recipientId,
	}
}

// newMessageEnvelope is the constructor for a messageEnvelope
func newMessageEnvelope(text string) *messageEnvelope {
	return &messageEnvelope{
		Text: text,
	}
}

// newTextToUserEnvelope is the constructor for a textToUserEnvelope
func newTextToUserEnvelope(recipientId, text string) *textToUserEnvelope {
	return &textToUserEnvelope{
		Recipient: &recipientEnvelope{
			Id: recipientId,
		},
		Message: &messageEnvelope{
			Text: text,
		},
	}
}
