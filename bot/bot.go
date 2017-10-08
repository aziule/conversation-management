package bot

import (
	"github.com/aziule/conversation-management/platform"
	"github.com/aziule/conversation-management/conversation"
)

// The main Bot structure
type bot struct {
	uuid string
	platform platform.Platform
	entrypoint *conversation.Entrypoint
	webhooks []*Webhook
}

// Getters
func (bot *bot) Uuid() string { return bot.uuid }
func (bot *bot) Platform() platform.Platform { return bot.platform }
func (bot *bot) Entrypoint() *conversation.Entrypoint { return bot.entrypoint }
func (bot *bot) Webhooks() []*Webhook { return bot.webhooks }