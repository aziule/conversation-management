package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/app/facebook/api"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
)

// Config is the config required in order to instantiate a new FacebookBot
type Config struct {
	VerifyToken            string
	ApiVersion             string
	PageAccessToken        string
	NlpParser              nlp.Parser
	ConversationRepository conversation.Repository
}

// Bot is the main structure
type facebookBot struct {
	webhooks               []*bot.Webhook
	verifyToken            string
	fbApi                  *api.FacebookApi
	nlpParser              nlp.Parser
	conversationRepository conversation.Repository
}

// NewBot is the constructor method that creates a Facebook bot, using
// the Config struct as method parameters.
// By default, we attach webhooks to the bot when constructing it.
// Later on, we can think about managing webhooks as we would manage events, and
// subscribe to the ones we like (for example, as defined in the conf).
func NewBot(config *Config) bot.Bot {
	bot := &facebookBot{
		verifyToken:            config.VerifyToken,
		fbApi:                  api.NewFacebookApi(config.ApiVersion, config.PageAccessToken, http.DefaultClient),
		nlpParser:              config.NlpParser,
		conversationRepository: config.ConversationRepository,
	}

	bot.bindDefaultWebhooks()

	return bot
}

// bindDefaultWebhooks initialises the default Facebook-related webhooks
// Use this method to create and bind the default Facebook webhooks to the bot
func (facebookBot *facebookBot) bindDefaultWebhooks() {
	facebookBot.webhooks = append(facebookBot.webhooks, bot.NewWebHook(
		bot.HttpMethodGet,
		"/",
		facebookBot.HandleValidateWebhook,
	))

	facebookBot.webhooks = append(facebookBot.webhooks, bot.NewWebHook(
		bot.HttpMethodPost,
		"/",
		facebookBot.HandleMessageReceived,
	))
}

// Webhooks returns the bot's webhooks.
// This method is required in order to inherit from the Bot interface.
func (facebookBot *facebookBot) Webhooks() []*bot.Webhook {
	return facebookBot.webhooks
}
