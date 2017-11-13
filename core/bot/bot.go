// Package bot defines the generic structs and interfaces for creating a bot
// on any platform.
package bot

// Bot is the main interface for a Bot
type Bot interface {
	Webhooks() []*Webhook
	ApiEndpoints() []*ApiEndpoint
}
