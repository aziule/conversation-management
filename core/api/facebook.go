package api

import (
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrFacebookApiNotFound = errors.New("Facebook API not found")

	// facebookApiBuilders stores the available FacebookApi builders
	facebookApiBuilders = make(map[string]FacebookApiBuilder)
)

// FacebookApiBuilder is the interface describing a builder for FacebookApi
type FacebookApiBuilder func(conf map[string]interface{}) (FacebookApi, error)

// RegisterFacebookApiBuilder adds a new FacebookApiBuilder to the list of available builders
func RegisterFacebookApiBuilder(name string, builder FacebookApiBuilder) {
	_, registered := facebookApiBuilders[name]

	if registered {
		log.WithField("name", name).Warning("FacebookApiBuilder already registered, ignoring")
	}

	facebookApiBuilders[name] = builder
}

// NewFacebookApi tries to create a FacebookApi using the available builders.
// Returns ErrFacebookApiNotFound if the facebookApi builder isn't found.
// Returns an error in case of any error during the build process.
func NewFacebookApi(name string, conf map[string]interface{}) (FacebookApi, error) {
	facebookApiBuilder, ok := facebookApiBuilders[name]

	if !ok {
		return nil, ErrFacebookApiNotFound
	}

	facebookApi, err := facebookApiBuilder(conf)

	if err != nil {
		return nil, err
	}

	return facebookApi, nil
}

// FacebookApi is the interface representing a Facebook API
// @todo: rename to Api
type FacebookApi interface {
	ParseRequestMessageReceived(r *http.Request) (*ReceivedMessage, error)
	SendTextToUser(recipientId, text string) error
}

// ReceivedMessage is the base struct for received messages
// @todo: see how to rename to FacebookReceivedMessage if facebook.go
// is the only file in the api package
type ReceivedMessage struct {
	Mid               string
	SenderId          string
	RecipientId       string
	SentAt            time.Time
	Text              string
	QuickReplyPayload string
	Nlp               []byte
}
