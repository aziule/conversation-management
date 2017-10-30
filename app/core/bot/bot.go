package bot

import (
	"github.com/aziule/conversation-management/app/core/conversation"
	"github.com/aziule/conversation-management/app/core/nlp"
)

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	NlpParser() nlp.Parser
	ConversationRepository() conversation.Repository
}
