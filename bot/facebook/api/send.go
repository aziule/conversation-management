package api

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"bytes"
)

// SendTextToUser is the FacebookApi's interface method responsible for sending a 1-to-1 message to a user
func (api *FacebookApi) SendTextToUser(recipientId, text string) error {
	url := api.getSendTextUrl()

	object := newTextToUserEnvelope(recipientId, text)
	jsonObject, err := json.Marshal(object)

	if err != nil {
		// @todo: fix
		return err
	}

	request, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(jsonObject))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		// @todo: fix
		return err
	}

	client := http.DefaultClient
	response, err := client.Do(request)

	if err != nil {
		// @todo: fix
		return err
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

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
