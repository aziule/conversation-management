package facebook

import (
	"github.com/aziule/conversation-management/core/bot"
	"net/http"
	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/nlu"
)

// Bot is the main structure
type facebookBot struct {
	pageAccessToken string
	fbApi           *FacebookApi
	nluParser nlu.Parser
}

// NewFacebookBot is the constructor method that creates a Facebook bot
func NewFacebookBot(config *core.Config, nluParser nlu.Parser) bot.Bot {
	return &facebookBot{
		pageAccessToken: config.FbPageAccessToken,
		fbApi:           NewFacebookApi(config.FbApiVersion, config.FbPageAccessToken, http.DefaultClient),
		nluParser: nluParser,
	}
}

// Webhooks returns the available webhooks for the bot
func (facebookBot *facebookBot) Webhooks() []*bot.Webhook {
	webhooks := []*bot.Webhook{}

	webhooks = append(webhooks, bot.NewWebHook(
		bot.HTTP_METHOD_GET,
		"/",
		HandleValidateWebhook,
	))

	webhooks = append(webhooks, bot.NewWebHook(
		bot.HTTP_METHOD_POST,
		"/",
		HandleMessageReceived,
	))

	return webhooks
}

// Getters
func (facebookBot *facebookBot) NluParser() nlu.Parser { return facebookBot.nluParser }
