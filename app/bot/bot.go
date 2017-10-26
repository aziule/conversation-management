package bot

import (
	"github.com/aziule/conversation-management/app/conversation"
	"github.com/aziule/conversation-management/app/nlp"
)

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	NlpParser() nlp.Parser
	ConversationRepository() conversation.Repository
}
