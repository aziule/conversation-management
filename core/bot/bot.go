package bot

import "github.com/aziule/conversation-management/core/nlp"

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	NlpParser() *nlp.Parser
}
