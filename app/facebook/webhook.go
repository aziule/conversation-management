package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/utils"
	log "github.com/sirupsen/logrus"
)

// handleMessageReceived is called when a new message is sent by the user to the page
// We parse the message, extract relevant NLP data, check the context, validate
// the data, and return a response.
func (bot *facebookBot) handleMessageReceived(w http.ResponseWriter, r *http.Request) {
	log.Debug("New Facebook message received")

	receivedMessage, err := bot.fbApi.ParseRequestMessageReceived(r)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("message", receivedMessage).Errorf("Could not parse the received message: %s", err)
		return
	}

	user, err := bot.conversationHandler.GetUser(receivedMessage.SenderId)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("user", receivedMessage.SenderId).Errorf("Could not find the user: %s", err)
		return
	}

	c, err := bot.conversationHandler.GetConversation(user)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("user", user).Infof("Could not get the conversation: %s", err)
		return
	}

	log.WithField("conversation", c).Debug("Conversation fetched")

	userMessage := conversation.NewUserMessage(
		receivedMessage.Text,
		receivedMessage.SentAt,
		user,
		nil,
	)

	// @todo: use the conversation handler instead
	c.AddMessage(userMessage)

	if receivedMessage.Nlp == nil {
		// @todo: handle this case: parse the text using the NLP parser
		log.Errorf("No data to parse")
		return
	}

	parsedData, err := bot.nlpParser.ParseNlpData(receivedMessage.Nlp)

	if err != nil {
		// @todo: handle this case and return something to the user. Make sure the
		// conversation is saved with the message. For example, we could think
		// about adding a flag to the message, like:
		// - could_not_parse_nlp
		// - could_not_process
		// - something_else
		// - ...
		// => gives more context and allows us to save data & understand it even
		// though errors occur.
		// @todo: save the conversation
		log.WithField("nlp", receivedMessage.Nlp).Errorf("Could not parse NLP data: %s", err)
		return
	}

	userMessage.ParsedData = parsedData

	log.WithField("data", parsedData).Debug("Data parsed from message")

	// @todo: save the conversation

	err = bot.conversationHandler.ProcessData(parsedData, c)

	if err != nil {
		log.WithFields(log.Fields{
			"data":         parsedData,
			"conversation": c,
		}).Errorf("Could not process the data: %s", err)
	}
}

// handleValidateWebhook tries to validate the Facebook webhook
// More information here: https://developers.facebook.com/docs/messenger-platform/getting-started/quick-start
func (bot *facebookBot) handleValidateWebhook(w http.ResponseWriter, r *http.Request) {
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

// bindDefaultWebhooks initialises the default Facebook-related webhooks.
// Use this method to create and bind the default Facebook webhooks to the bot.
func (b *facebookBot) bindDefaultWebhooks() {
	b.webhooks = append(b.webhooks, bot.NewWebHook(
		"GET",
		"/fb",
		b.handleValidateWebhook,
	))

	b.webhooks = append(b.webhooks, bot.NewWebHook(
		"POST",
		"/fb",
		b.handleMessageReceived,
	))
}
