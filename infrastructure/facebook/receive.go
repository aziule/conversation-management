package facebook

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/aziule/conversation-management/core/api"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCouldNotReadRequestBody = errors.New("Could not read the request's body")
	ErrInvalidJson             = errors.New("Invalid JSON")
	ErrMissingKey              = func(key string) error { return errors.New(fmt.Sprintf("Missing key: %s", key)) }
	ErrNoEntry                 = errors.New("No entry to parse")
	ErrNoMessage               = errors.New("No message to parse")
)

// ParseJsonBody creates a Message from json bytes and returns an error if a parsing issue occurred
func (fbApi *facebookApi) ParseRequestMessageReceived(r *http.Request) (*api.FacebookReceivedMessage, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Infof("Could not read the request's body: %s", err)
		return nil, ErrCouldNotReadRequestBody
	}

	json, err := jason.NewObjectFromBytes(body)

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

	// Get the number of seconds
	sentAt = sentAt / 1000

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

	return &api.FacebookReceivedMessage{
		Mid:               mid,
		SenderId:          senderId,
		RecipientId:       recipientId,
		SentAt:            time.Unix(sentAt, 0),
		Text:              text,
		QuickReplyPayload: quickReplyPayload,
		Nlp:               nlpBytes,
	}, nil
}
