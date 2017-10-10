package bot

import (
	"github.com/aziule/conversation-management/conversation"
	"github.com/aziule/conversation-management/facebook"
)

// Bot is the main structure
type Bot struct {
	pageAccessToken string
	fbApi *facebook.FacebookApi
}

// NewBot is the constructor method
func NewBot() *Bot {
	return &Bot{
		pageAccessToken: "EAALQ00uMgf8BAHnz4y711RUBOQyiQrUpy4ZAIXeXvL4L0mIZAZCq6WXKZBnwwhT8Xfw2So5DZABaRSfxjuO97mdQklTxZCdZATKFH7xvJ5VEwqsCyQRTXh9yTq9ZBSGATaSZCSsS7xhv3TeHvFyx5s0xcQ88BxiZBqmv8zPFRcTX9iJAZDZD",
	}
}

// Getters
func (bot *Bot) PageAccessToken() string { return bot.pageAccessToken }

// SentText sends text to a given user
func (bot *Bot) sendText(text string, user *conversation.User) error {
	return nil
}

// LoadStories loads the base stories of the bot
func (bot *Bot) loadStories() *[]conversation.Story {
	return nil
}
