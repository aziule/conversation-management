package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/bot/facebook/api"
	"github.com/aziule/conversation-management/core"
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/nlp/wit"
)

// Bot is the main structure
type facebookBot struct {
	verifyToken string
	fbApi       *api.FacebookApi
	webhooks    []*bot.Webhook
	nlpParser   nlp.Parser
}

// NewFacebookBot is the constructor method that creates a Facebook bot
// By default, we attach webhooks to the bot when constructing it.
// Later on, we can think about managing webhooks as we would manage events, and
// subscribe to the ones we like (for example, as defined in the conf).
func NewFacebookBot(config *core.Config) *facebookBot {
	dataTypeMap := getDefaultDataTypeMap()

	bot := &facebookBot{
		verifyToken: config.FbVerifyToken,
		fbApi:       api.NewFacebookApi(config.FbApiVersion, config.FbPageAccessToken, http.DefaultClient),
		nlpParser:   wit.NewParser(dataTypeMap),
	}

	bot.BindDefaultWebhooks()

	return bot
}

// InitWebhooks initialises the default Facebook-related webhooks
// Use this method to create and bind the default Facebook webhooks to the bot
func (facebookBot *facebookBot) BindDefaultWebhooks() {
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

// getDefaultDataTypeMap returns the default data type map.
// For now, this is highly coupled with Wit's data types and should
// be updated every time a change is made to Wit.
func getDefaultDataTypeMap() nlp.DataTypeMap {
	DataTypeMap := make(nlp.DataTypeMap)

	DataTypeMap["nb_persons"] = nlp.DataTypeInt
	DataTypeMap["intent"] = nlp.DataTypeIntent
	DataTypeMap["datetime"] = nlp.DataTypeDateTime

	return DataTypeMap
}

// Getters
func (facebookBot *facebookBot) Webhooks() []*bot.Webhook { return facebookBot.webhooks }
func (facebookBot *facebookBot) NlpParser() nlp.Parser    { return facebookBot.nlpParser }
