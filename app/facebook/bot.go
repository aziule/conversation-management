// Package facebook defines Facebook-related bot methods and behaviour.
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
	StoryRepository        conversation.StoryRepository
}

// facebookBot is the main structure
type facebookBot struct {
	webhooks            []*bot.Webhook
	verifyToken         string
	fbApi               *api.FacebookApi
	nlpParser           nlp.Parser
	conversationHandler conversation.Handler
}

// NewBot is the constructor method that creates a Facebook bot, using
// the Config struct as method parameters.
//
// Upon creation:
// - The webhooks are attached.
// - We load the list of stories.
func NewBot(config *Config) *facebookBot {
	bot := &facebookBot{
		verifyToken: config.VerifyToken,
		fbApi:       api.NewFacebookApi(config.ApiVersion, config.PageAccessToken, http.DefaultClient),
		nlpParser:   config.NlpParser,
	}

	// @todo: we need to check if all of the stories's steps are being handled
	// and if any are missing / deprecated.
	bot.conversationHandler = newConversationHandler(
		bot.getDefaultStepsMapping(), // @todo: directly pass the step handler rather than the steps mapping
		config.ConversationRepository,
		config.StoryRepository,
	)

	bot.bindDefaultWebhooks()

	return bot
}

// Webhooks returns the bot's webhooks.
// This method is required in order to inherit from the Bot interface.
func (b *facebookBot) Webhooks() []*bot.Webhook {
	return b.webhooks
}

// bindDefaultWebhooks initialises the default Facebook-related webhooks.
// Use this method to create and bind the default Facebook webhooks to the bot.
func (b *facebookBot) bindDefaultWebhooks() {
	b.webhooks = append(b.webhooks, bot.NewWebHook(
		bot.HttpMethodGet,
		"/",
		b.handleValidateWebhook,
	))

	b.webhooks = append(b.webhooks, bot.NewWebHook(
		bot.HttpMethodPost,
		"/",
		b.handleMessageReceived,
	))
}

// getDefaultStepsMapping returns the default steps mapping between
// a step's name and its handling func.
// @todo: find a better name and/or move somewhere else
func (b *facebookBot) getDefaultStepsMapping() conversation.StepsProcessMap {
	pm := conversation.StepsProcessMap{}

	pm["get_intent"] = b.processStepGetIntent

	return pm
}
