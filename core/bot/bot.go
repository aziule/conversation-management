package bot

import (
	"github.com/aziule/conversation-management/core/nlu"
)

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	NluParser() nlu.Parser
}
