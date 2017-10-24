package conversation

import (
	"github.com/aziule/conversation-management/core/nlp"
	"time"
)

type Message interface {
	Text() string
	SentAt() time.Time
}

type message struct {
	Text   string
	SentAt time.Time
}

type UserMessage struct {
	message
	Sender     User
	ParsedData nlp.ParsedData
}

type BotMessage struct {
	message
	RepliesTo *UserMessage
}

type MessagesFlow struct {
	Messages []Message
}

func (flow *MessagesFlow) IsNew() bool {
	return len(flow.Messages) == 0
}

func (flow *MessagesFlow) LastMessage() Message {
	if len(flow.Messages) == 0 {
		return nil
	}

	var lastMessage Message

	for _, message := range flow.Messages {
		if message.SentAt() > lastMessage.SentAt() {
			lastMessage = message
		}
	}

	return lastMessage
}
