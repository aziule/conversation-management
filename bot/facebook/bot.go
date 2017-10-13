package facebook

import (
	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/nlu"
	"net/http"
)

// Bot is the main structure
type facebookBot struct {
	pageAccessToken string
	verifyToken     string
	fbApi           *FacebookApi
	nluParser       nlu.Parser
}

// NewFacebookBot is the constructor method that creates a Facebook bot
func NewFacebookBot(config *core.Config) bot.Bot {
	return &facebookBot{
		pageAccessToken: config.FbPageAccessToken,
		verifyToken:     config.FbVerifyToken,
		fbApi:           NewFacebookApi(config.FbApiVersion, config.FbPageAccessToken, http.DefaultClient),
		nluParser:       nlu.NewParserFromConfig(config),
	}
}

// Webhooks returns the available webhooks for the bot
func (facebookBot *facebookBot) Webhooks() []*bot.Webhook {
	webhooks := []*bot.Webhook{}

	webhooks = append(webhooks, bot.NewWebHook(
		bot.HTTP_METHOD_GET,
		"/",
		facebookBot.HandleMessageReceived,
	))

	webhooks = append(webhooks, bot.NewWebHook(
		bot.HTTP_METHOD_POST,
		"/",
		facebookBot.HandleMessageReceived,
	))

	return webhooks
}

// Getters
func (facebookBot *facebookBot) NluParser() nlu.Parser   { return facebookBot.nluParser }
func (facebookBot *facebookBot) PageAccessToken() string { return facebookBot.pageAccessToken }
func (facebookBot *facebookBot) VerifyToken() string     { return facebookBot.verifyToken }
