package api

import (
	"net/http"
	"time"

	"github.com/aziule/conversation-management/core/utils"
)

const builderPrefix = "api_"

// RegisterFacebookApiBuilder registers a new service builder
func RegisterFacebookApiBuilder(name string, builder utils.ServiceBuilder) {
	utils.RegisterServiceBuilder(builderPrefix+name, builder)
}

// NewFacebookApi tries to create a FacebookApi using the available builders.
// Returns ErrFacebookApiNotFound if the facebookApi builder isn't found.
// Returns an error in case of any error during the build process.
func NewFacebookApi(name string, conf utils.BuilderConf) (FacebookApi, error) {
	facebookApiBuilder, err := utils.GetServiceBuilder(builderPrefix + name)

	if err != nil {
		return nil, err
	}

	facebookApi, err := facebookApiBuilder(conf)

	if err != nil {
		return nil, err
	}

	return facebookApi.(FacebookApi), nil
}

// FacebookApi is the interface representing a Facebook API
type FacebookApi interface {
	ParseRequestMessageReceived(r *http.Request) (*FacebookReceivedMessage, error)
	SendTextToUser(recipientId, text string) error
}

// FacebookReceivedMessage is the base struct for received messages
// @todo: see how to rename to FacebookFacebookReceivedMessage if facebook.go
// is the only file in the api package
type FacebookReceivedMessage struct {
	Mid               string
	SenderId          string
	RecipientId       string
	SentAt            time.Time
	Text              string
	QuickReplyPayload string
	Nlp               []byte
}
