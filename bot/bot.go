package bot

import (
	"github.com/aziule/conversation-management/conversation"
	"github.com/aziule/conversation-management/nlp"
)

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	NlpParser() nlp.Parser
	ConversationReader() conversation.Reader
	ConversationWriter() conversation.Writer
}
