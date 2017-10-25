package facebook

import (
	"net/http"

	"github.com/aziule/conversation-management/bot"
	"github.com/aziule/conversation-management/bot/facebook/api"
	"github.com/aziule/conversation-management/conversation"
	"github.com/aziule/conversation-management/nlp"
)

// DefaultDataTypeMap is the default data type map to be used with Wit.
// For now, this is highly coupled with Wit's data types and should
// be updated every time a change is made to Wit.
// It is initialised in the init() function
var DefaultDataTypeMap nlp.DataTypeMap

// Config is the config required in order to instantiate a new FacebookBot
type Config struct {
	VerifyToken        string
	ApiVersion         string
	PageAccessToken    string
	NlpParser          nlp.Parser
	ConversationReader conversation.Reader
	ConversationWriter conversation.Writer
}

// Bot is the main structure
type facebookBot struct {
	verifyToken        string
	fbApi              *api.FacebookApi
	webhooks           []*bot.Webhook
	nlpParser          nlp.Parser
	conversationReader conversation.Reader
	conversationWriter conversation.Writer
}

// NewFacebookBot is the constructor method that creates a Facebook bot, using
// the Config struct as method parameters.
// By default, we attach webhooks to the bot when constructing it.
// Later on, we can think about managing webhooks as we would manage events, and
// subscribe to the ones we like (for example, as defined in the conf).
func NewBot(config *Config) bot.Bot {
	bot := &facebookBot{
		verifyToken:        config.VerifyToken,
		fbApi:              api.NewFacebookApi(config.ApiVersion, config.PageAccessToken, http.DefaultClient),
		nlpParser:          config.NlpParser,
		conversationReader: config.ConversationReader,
		conversationWriter: config.ConversationWriter,
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

// Webhooks returns the bot's webhooks.
// This method is required in order to inherit from the Bot interface.
func (facebookBot *facebookBot) Webhooks() []*bot.Webhook {
	return facebookBot.webhooks
}

// NlpParser returns the bot's nlp parser.
// This method is required in order to inherit from the Bot interface.
func (facebookBot *facebookBot) NlpParser() nlp.Parser {
	return facebookBot.nlpParser
}

// ConversationReader returns the ConversationReader.
// This method is required in order to inherit from the Bot interface.
func (facebookBot *facebookBot) ConversationReader() conversation.Reader {
	return facebookBot.conversationReader
}

// ConversationWriter returns the ConversationWriter.
// This method is required in order to inherit from the Bot interface.
func (facebookBot *facebookBot) ConversationWriter() conversation.Writer {
	return facebookBot.conversationWriter
}

func init() {
	DefaultDataTypeMap = make(nlp.DataTypeMap)
	DefaultDataTypeMap["nb_persons"] = nlp.DataTypeInt
	DefaultDataTypeMap["intent"] = nlp.DataTypeIntent
	DefaultDataTypeMap["datetime"] = nlp.DataTypeDateTime
}
