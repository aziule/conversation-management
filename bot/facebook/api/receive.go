package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/antonholmquist/jason"
)

// ReceivedMessage is the base struct for received messages
type ReceivedMessage struct {
	mid               string
	senderId          string
	recipientId       string
	sentAt            time.Time
	text              string
	quickReplyPayload string
	nlp               []byte
}

// ParseJsonBody creates a Message from json bytes and returns an error if a parsing issue occurred
func (api *FacebookApi) ParseRequestMessageReceived(r *http.Request) (*ReceivedMessage, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		// @todo: log and handle error types
		return nil, err
	}

	json, err := jason.NewObjectFromBytes(body)

	prettyPrint(body)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Could not parse JSON") // @todo: error types
	}

	// The message content itself is embedded inside the "entry" array
	entries, err := json.GetObjectArray("entry")

	if err != nil {
		return nil, errors.New("Missing entry key")
	}

	if len(entries) == 0 {
		return nil, errors.New("No message to parse")
	}

	entry := entries[0]

	messaging, err := entry.GetObjectArray("messaging")

	if err != nil {
		return nil, errors.New("Missing messaging key")
	}

	if len(messaging) == 0 {
		return nil, errors.New("No message to parse")
	}

	messageData := messaging[0]

	mid, err := messageData.GetString("message", "mid")

	if err != nil {
		return nil, errors.New("Missing message.id")
	}

	senderId, err := messageData.GetString("sender", "id")

	if err != nil {
		return nil, errors.New("Missing sender.id")
	}

	recipientId, err := messageData.GetString("recipient", "id")

	if err != nil {
		return nil, errors.New("Missing recipient.id")
	}

	sentAt, err := messageData.GetInt64("timestamp")

	if err != nil {
		return nil, errors.New("Missing timestamp")
	}

	text, _ := messageData.GetString("message", "text")
	quickReplyPayload, _ := messageData.GetString("quick_reply", "payload")

	nlp, err := messageData.GetObject("message", "nlp", "entities")
	var nlpBytes []byte

	if err == nil {
		nlpBytes, err = nlp.MarshalJSON()

		if err != nil {
			// @todo: log error
		}
	} else {
		// @todo: log that no NLP was received
	}

	return &ReceivedMessage{
		mid:               mid,
		senderId:          senderId,
		recipientId:       recipientId,
		sentAt:            time.Unix(sentAt, 0),
		text:              text,
		quickReplyPayload: quickReplyPayload,
		nlp:               nlpBytes,
	}, nil
}

// Getters
func (m *ReceivedMessage) SenderId() string          { return m.senderId }
func (m *ReceivedMessage) RecipientId() string       { return m.recipientId }
func (m *ReceivedMessage) SentAt() time.Time         { return m.sentAt }
func (m *ReceivedMessage) Mid() string               { return m.mid }
func (m *ReceivedMessage) Text() string              { return m.text }
func (m *ReceivedMessage) QuickReplyPayload() string { return m.quickReplyPayload }
func (m *ReceivedMessage) Nlp() []byte               { return m.nlp }
