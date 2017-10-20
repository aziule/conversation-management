package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"fmt"
	"github.com/antonholmquist/jason"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCouldNotReadRequestBody = errors.New("Could not read the request's body")
	ErrInvalidJson             = errors.New("Invalid JSON")
	ErrMissingKey              = func(key string) error { return errors.New(fmt.Sprintf("Missing key: %s", key)) }
	ErrNoEntry                 = errors.New("No entry to parse")
	ErrNoMessage               = errors.New("No message to parse")
)

// ReceivedMessage is the base struct for received messages
type ReceivedMessage struct {
	Mid               string
	SenderId          string
	RecipientId       string
	SentAt            time.Time
	Text              string
	QuickReplyPayload string
	Nlp               []byte
}

// ParseJsonBody creates a Message from json bytes and returns an error if a parsing issue occurred
func (api *FacebookApi) ParseRequestMessageReceived(r *http.Request) (*ReceivedMessage, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Infof("Could not read the request's body: %s", err)
		return nil, ErrCouldNotReadRequestBody
	}

	json, err := jason.NewObjectFromBytes(body)

	// @todo: remove (use with debug only)
	prettyPrint(body)

	if err != nil {
		log.WithField("body", string(body)).Infof("Could not parse JSON from the request: %s", err)
		return nil, ErrInvalidJson
	}

	// The message content itself is embedded inside the "entry" array
	entries, err := json.GetObjectArray("entry")

	if err != nil {
		log.WithField("key", "entry").Info("Missing key")
		return nil, ErrMissingKey("entry")
	}

	if len(entries) == 0 {
		log.Info("No entries to parse")
		return nil, ErrNoEntry
	}

	entry := entries[0]

	messaging, err := entry.GetObjectArray("messaging")

	if err != nil {
		log.WithField("key", "messaging").Info("Missing key")
		return nil, ErrMissingKey("messaging")
	}

	if len(messaging) == 0 {
		log.Info("No message to parse")
		return nil, ErrNoMessage
	}

	messageData := messaging[0]

	mid, err := messageData.GetString("message", "mid")

	if err != nil {
		log.WithField("key", "message.id").Info("Missing key")
		return nil, ErrMissingKey("message.id")
	}

	senderId, err := messageData.GetString("sender", "id")

	if err != nil {
		log.WithField("key", "sender.id").Info("Missing key")
		return nil, ErrMissingKey("sender.id")
	}

	recipientId, err := messageData.GetString("recipient", "id")

	if err != nil {
		log.WithField("key", "recipient.id").Info("Missing key")
		return nil, ErrMissingKey("recipient.id")
	}

	sentAt, err := messageData.GetInt64("timestamp")

	if err != nil {
		log.WithField("key", "timestamp").Info("Missing key")
		return nil, ErrMissingKey("timestamp")
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
		Mid:               mid,
		SenderId:          senderId,
		RecipientId:       recipientId,
		SentAt:            time.Unix(sentAt, 0),
		Text:              text,
		QuickReplyPayload: quickReplyPayload,
		Nlp:               nlpBytes,
	}, nil
}
