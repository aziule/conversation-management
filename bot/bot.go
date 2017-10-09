package bot

import "github.com/aziule/conversation-management/conversation"

// Bot is the main structure
type Bot struct {
	version float32
}

// Getters
func (bot *Bot) Version() int { return bot.version }

// LoadStories loads the base stories of the bot
func (bot *Bot) LoadStories() *[]conversation.Story {
	return nil
}