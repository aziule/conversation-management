package bot

// Bot is the main interface for a Bot
type Bot interface {
	BindWebhooks([]*Webhook)
	Webhooks() []*Webhook
}
