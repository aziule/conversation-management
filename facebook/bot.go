package facebook

import (
	"github.com/aziule/conversation-management/core"
	"net/http"
)

// Bot is the main structure
type facebookBot struct {
	pageAccessToken string
	fbApi           *FacebookApi
}

// NewFacebookBot is the constructor method that creates a Facebook bot
func NewFacebookBot(config *core.Config) core.Bot {
	return &facebookBot{
		pageAccessToken: config.FbPageAccessToken,
		fbApi:           NewFacebookApi(config.FbApiVersion, config.FbPageAccessToken, http.DefaultClient),
	}
}

// Webhooks returns the available webhooks for the bot
func (facebookBot *facebookBot) Webhooks() []*core.Webhook {
	webhooks := []*core.Webhook{}

	webhooks = append(webhooks, core.NewWebHook(
		core.HTTP_METHOD_GET,
		"/",
		HandleValidateWebhook,
	))

	webhooks = append(webhooks, core.NewWebHook(
		core.HTTP_METHOD_POST,
		"/",
		HandleMessageReceived,
	))

	return webhooks
}
