package conversation

import (
	"time"

	"github.com/aziule/conversation-management/core/nlp"
)

type Message interface {
	Channel() Channel
	Text() string
	SentAt() time.Time
}

type UserMessage interface {
	Message
	Sender() User
	ParsedData() nlp.ParsedData
}

type BotMessage interface {
	Message
	RepliesTo() UserMessage
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
