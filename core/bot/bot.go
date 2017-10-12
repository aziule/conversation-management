package bot

import (
	"github.com/aziule/conversation-management/core/nlu"
	"github.com/aziule/conversation-management/core"
)

// Platform represents the available platforms for a given Bot
type Platform string

const (
	PLATFORM_FACEBOOK Platform = "facebook" // Is a bot on Facebook
)

var (
	botFactories = make(map[Platform]BotFactory) // The list of available factories
)

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	NluParser() nlu.Parser
}

// BotFactory is the main factory func used to instantiate new Bot implementations
type BotFactory func(*core.Config) Bot

// RegisterFactory allows us to register factory methods for creating new bots
func RegisterFactory(platform Platform, factory BotFactory) {
	if factory == nil {
		return
	}

	_, registered := botFactories[platform]

	if registered {
		return
	}

	botFactories[platform] = factory
}

// NewParserFromConfig instantiates the correct Bot given the platform and the configuration
func NewBotFromConfig(platform Platform, config *core.Config) Bot {
	return botFactories[platform](config)
}
