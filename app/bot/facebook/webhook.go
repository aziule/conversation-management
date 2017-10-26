package facebook

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aziule/conversation-management/app/conversation"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

var ErrCouldNotFetchParam = func(key string) error { return errors.New(fmt.Sprintf("Could not fetch param: %s", key)) }

// HandleMessageReceived is called when a new message is sent by the user to the page
// We parse the message, extract relevant NLP data, check the context, validate
// the data, and return a response.
func (bot *facebookBot) HandleMessageReceived(w http.ResponseWriter, r *http.Request) {
	log.Debug("New Facebook message received")

	receivedMessage, err := bot.fbApi.ParseRequestMessageReceived(r)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.Errorf("Could not parse the received message: %s", err.Error())
		return
	}

	if receivedMessage.Nlp == nil {
		// @todo: handle this case
		log.Errorf("No data to parse")
		return
	}

	parsedData, err := bot.nlpParser.ParseNlpData(receivedMessage.Nlp)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.Errorf("Could not parse NLP data: %s", err.Error())
		return
	}

	user, err := bot.conversationRepository.FindUser(receivedMessage.SenderId)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.Errorf("Could not find the user: %s", receivedMessage.SenderId)
		return
	}

	if user == nil {
		log.Debugf("Inserting the user: %s", receivedMessage.SenderId)

		// Insert the user
		err = bot.conversationRepository.InsertUser(
			&conversation.User{
				Id:   bson.NewObjectId(),
				FbId: receivedMessage.SenderId,
			},
		)

		if err != nil {
			// @todo: handle this case and return something to the user
			log.Errorf("Could not insert the user: %s", receivedMessage.SenderId)
			return
		}
	}

	log.Debugf("Request from user: %s", receivedMessage.SenderId)
	c, err := bot.conversationRepository.FindLatestConversation(user)

	if err != nil {
		if err != conversation.ErrNotFound {
			// @todo: handle this case and return something to the user
			log.Errorf("Could not find the latest conversation: %s", receivedMessage.SenderId)
			return
		}

		// The conversation was not found: create a new one
		c = conversation.NewConversation()
	}

	// Create a new conversation if the previous one is over
	if c.Status == conversation.StatusOver {
		c = conversation.NewConversation()
	}

	userMessage := conversation.NewUserMessage(
		receivedMessage.Text,
		receivedMessage.SentAt,
		user,
		parsedData,
	)

	log.Debug(userMessage)
	//c.Received(userMessage)

	bot.fbApi.SendTextToUser(receivedMessage.SenderId, receivedMessage.Text)
}

// HandleValidateWebhook tries to validate the Facebook webhook
// More information here: https://developers.facebook.com/docs/messenger-platform/getting-started/quick-start
func (bot *facebookBot) HandleValidateWebhook(w http.ResponseWriter, r *http.Request) {
	log.Debug("New Facebook webhook validation request")

	queryParams := r.URL.Query()

	hubMode, err := getSingleQueryParam(queryParams, "hub.mode")

	if err != nil {
		log.Infof("Could not fetch param: %s", err)
		return
	}

	if hubMode != "subscribe" {
		log.Debugf("Invalid hub mode: %s", hubMode)
		return
	}

	verifyToken, err := getSingleQueryParam(queryParams, "hub.verify_token")

	if err != nil {
		log.Infof("Could not fetch param: %s", err)
		return
	}

	if verifyToken != bot.verifyToken {
		log.Debugf("Invalid verify token: %s", verifyToken)
		return
	}

	challenge, err := getSingleQueryParam(queryParams, "hub.challenge")

	if err != nil {
		log.Infof("Could not fetch param: %s", err)
		return
	}

	// Validate the webhook by writing back the "hub.challenge" query param
	w.Write([]byte(challenge))
}

// getSingleQueryParam fetches a single query param using the given url values
func getSingleQueryParam(values url.Values, key string) (string, error) {
	params, ok := values[key]

	if !ok || len(params) != 1 {
		return "", ErrCouldNotFetchParam(key)
	}

	return params[0], nil
}
