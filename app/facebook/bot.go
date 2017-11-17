// Package facebook defines Facebook-related bot methods and behaviour.
package facebook

import (
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/infrastructure/facebook/api"
)

// Config is the config required in order to instantiate a new FacebookBot
type Config struct {
	VerifyToken            string
	FbApi                  api.FacebookApi
	NlpParser              nlp.Parser
	ConversationRepository conversation.Repository
	StoryRepository        conversation.StoryRepository
}

// facebookBot is the main structure
type facebookBot struct {
	verifyToken         string
	webhooks            []*bot.Webhook
	apiEndpoints        []*bot.ApiEndpoint
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
	}

	// @todo: we need to check if all of the stories's steps are being handled
	// and if any are missing / deprecated.
	bot.conversationHandler = newConversationHandler(
		bot.getDefaultStepsMapping(), // @todo: directly pass the step handler rather than the steps mapping
		config.ConversationRepository,
		config.StoryRepository,
		config.NlpParser,
		config.FbApi,
	)

	bot.bindDefaultWebhooks()
	bot.bindDefaultApiEndpoints()

	return bot
}

// Webhooks returns the bot's webhooks.
// This method is required in order to inherit from the Bot interface.
func (b *facebookBot) Webhooks() []*bot.Webhook {
	return b.webhooks
}

// ApiEndpoints returns the bot's available API endpoints.
// This method is required in order to inherit from the Bot interface.
func (b *facebookBot) ApiEndpoints() []*bot.ApiEndpoint {
	return b.apiEndpoints
}

// getDefaultStepsMapping returns the default steps mapping between
// a step's name and its handling func.
// @todo: find a better name and/or move somewhere else
func (b *facebookBot) getDefaultStepsMapping() conversation.StepsProcessMap {
	pm := conversation.StepsProcessMap{}

	pm["book_table_entrypoint"] = b.processStepBookTable
	pm["book_table_get_nb_persons"] = b.processStepBookTableGetNbPersons
	pm["book_table_get_time"] = b.processStepBookTableGetTime

	return pm
}
