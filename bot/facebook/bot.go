package facebook

import (
	"github.com/aziule/conversation-management/bot/facebook/api"
	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/bot"
	"net/http"
)

// Bot is the main structure
type facebookBot struct {
	pageAccessToken string
	verifyToken     string
	fbApi           *api.FacebookApi
	webhooks        []*bot.Webhook
}

// NewFacebookBot is the constructor method that creates a Facebook bot
// By default, no webhook is attached to the bot. It must be added manually
// using either the BindWebhooks or BindDefaultWebhooks method
func NewFacebookBot(config *core.Config) *facebookBot {
	return &facebookBot{
		pageAccessToken: config.FbPageAccessToken,
		verifyToken:     config.FbVerifyToken,
		fbApi:           api.NewFacebookApi(config.FbApiVersion, config.FbPageAccessToken, http.DefaultClient),
		webhooks:        nil,
	}
}

// Webhooks returns the available webhooks for the bot
func (facebookBot *facebookBot) Webhooks() []*bot.Webhook {
	return facebookBot.webhooks
}

// BindWebhooks binds the given webhooks to the bot
func (facebookBot *facebookBot) BindWebhooks(webhooks []*bot.Webhook) {
	for _, webhook := range webhooks {
		facebookBot.webhooks = append(facebookBot.webhooks, webhook)
	}
}

// InitWebhooks initialises the default Facebook-related webhooks
// Use this method to create and bind the default Facebook webhooks to the bot
func (facebookBot *facebookBot) BindDefaultWebhooks() {
	webhooks := []*bot.Webhook{}

	webhooks = append(webhooks, bot.NewWebHook(
		bot.HttpMethodGet,
		"/",
		facebookBot.HandleValidateWebhook,
	))

	webhooks = append(webhooks, bot.NewWebHook(
		bot.HttpMethodPost,
		"/",
		facebookBot.HandleMessageReceived,
	))

	facebookBot.BindWebhooks(webhooks)
}

// Getters
func (facebookBot *facebookBot) PageAccessToken() string { return facebookBot.pageAccessToken }
func (facebookBot *facebookBot) VerifyToken() string     { return facebookBot.verifyToken }
