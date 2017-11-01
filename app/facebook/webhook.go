package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/core/utils"
	log "github.com/sirupsen/logrus"
)

// HandleMessageReceived is called when a new message is sent by the user to the page
// We parse the message, extract relevant NLP data, check the context, validate
// the data, and return a response.
func (bot *facebookBot) HandleMessageReceived(w http.ResponseWriter, r *http.Request) {
	log.Debug("New Facebook message received")

	receivedMessage, err := bot.fbApi.ParseRequestMessageReceived(r)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.Errorf("Could not parse the received message: %s", err)
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
		log.WithField("nlp", receivedMessage.Nlp).Errorf("Could not parse NLP data: %s", err)
		return
	}

	log.WithField("data", parsedData).Debug("Data parsed from message")

	user, err := bot.conversationHandler.GetUser(receivedMessage.SenderId)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("user", receivedMessage.SenderId).Errorf("Could not find the user: %s", err)
		return
	}

	conversation, err := bot.conversationHandler.GetConversation(user)

	if err != nil {
		log.WithField("user", user).Infof("Could not get the conversation: %s", err)
	}

	log.WithField("conversation", conversation).Debug("Conversation fetched")
}

// HandleValidateWebhook tries to validate the Facebook webhook
// More information here: https://developers.facebook.com/docs/messenger-platform/getting-started/quick-start
func (bot *facebookBot) HandleValidateWebhook(w http.ResponseWriter, r *http.Request) {
	log.Debug("New Facebook webhook validation request")

	queryParams := r.URL.Query()

	hubMode, err := utils.GetSingleQueryParam(queryParams, "hub.mode")

	if err != nil {
		log.WithField("param", "hub.mode").Infof("Could not fetch param: %s", err)
		return
	}

	if hubMode != "subscribe" {
		log.WithFields(log.Fields{
			"expected": "subscribe",
			"mode":     hubMode,
		}).Info("Invalid hub mode")
		return
	}

	verifyToken, err := utils.GetSingleQueryParam(queryParams, "hub.verify_token")

	if err != nil {
		log.WithField("param", "hub.verify_token").Infof("Could not fetch param: %s", err)
		return
	}

	if verifyToken != bot.verifyToken {
		log.WithFields(log.Fields{
			"expected": bot.verifyToken,
			"token":    verifyToken,
		}).Info("Invalid verify token")
		return
	}

	challenge, err := utils.GetSingleQueryParam(queryParams, "hub.challenge")

	if err != nil {
		log.WithField("param", "hub.challenge").Infof("Could not fetch param: %s", err)
		return
	}

	// Validate the webhook by writing back the "hub.challenge" query param
	w.Write([]byte(challenge))
}
