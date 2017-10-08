package bot

import (
	"github.com/aziule/conversation-management/platform"
	"github.com/aziule/conversation-management/conversation"
)

// The main Bot structure
type Bot struct {
	uuid string
	entrypoint *conversation.Entrypoint
	webhooks []*Webhook
}

// Getters
func (bot *Bot) Uuid() string { return bot.uuid }
func (bot *Bot) Entrypoint() *conversation.Entrypoint { return bot.entrypoint }
func (bot *Bot) Webhooks() []*Webhook { return bot.webhooks }