package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/utils"
	log "github.com/sirupsen/logrus"
)

// handleMessageReceived is called when a new message is sent by the user to the page
// We parse the message, extract relevant NLP data, check the context, validate
// the data, and return a response.
func (bot *facebookBot) handleMessageReceived(w http.ResponseWriter, r *http.Request) {
	log.Debug("New Facebook message received")

	err := bot.conversationHandler.MessageReceived(r)

	if err != nil {
		log.Error(err)
	}

	log.Info("All good")
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
