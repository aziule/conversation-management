package core

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
}