// Package facebook defines Facebook-related bot methods and behaviour.
package facebook

import (
	"github.com/aziule/conversation-management/core/bot"
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/infrastructure/facebook/api"
)

const VerifyToken bot.ParamName = "verify_token"

// Config is the config required in order to instantiate a new FacebookBot
type Config struct {
	Metadata               *bot.Metadata
	FbApi                  api.FacebookApi
	NlpParser              nlp.Parser
	ConversationRepository conversation.Repository
	StoryRepository        conversation.StoryRepository
}

// facebookBot is the main structure
type facebookBot struct {
	webhooks            []*bot.Webhook
	apiEndpoints        []*bot.ApiEndpoint
	metadata            *bot.Metadata
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
		metadata: config.Metadata,
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
// This method is required in order to implement the Bot interface.
func (b *facebookBot) Webhooks() []*bot.Webhook {
	return b.webhooks
}

// ApiEndpoints returns the bot's available API endpoints.
// This method is required in order to implement the Bot interface.
func (b *facebookBot) ApiEndpoints() []*bot.ApiEndpoint {
	return b.apiEndpoints
}

// Metadata returns the bot's metadata.
// This method is required in order to implement the Bot interface.
func (b *facebookBot) Metadata() *bot.Metadata {
	return b.metadata
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
