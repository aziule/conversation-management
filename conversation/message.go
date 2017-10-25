package conversation

import (
	"github.com/aziule/conversation-management/nlp"
	"time"
)

type Message interface {
	Text() string
	SentAt() time.Time
}

// message is the struct that contains the base information about a generic message
type message struct {
	Text   string
	SentAt time.Time
}

// newMessage is the private constructor method for message
func newMessage(text string, sentAt time.Time) *message {
	return &message{
		Text:   text,
		SentAt: sentAt,
	}
}

// UserMessage represents a message received from a user
type UserMessage struct {
	*message
	Sender     *User
	ParsedData *nlp.ParsedData
}

// NewUserMessage is the constructor method for UserMessage
func NewUserMessage(text string, sentAt time.Time, sender *User, parsedData *nlp.ParsedData) *UserMessage {
	return &UserMessage{
		newMessage(text, sentAt),
		sender,
		parsedData,
	}
}

func (msg *UserMessage) Text() string {
	return msg.message.Text
}

func (msg *UserMessage) SentAt() time.Time {
	return msg.message.SentAt
}

type BotMessage struct {
	message
	RepliesTo *UserMessage
}

type messagesFlow struct {
	Messages []Message
}

func (flow *messagesFlow) IsNew() bool {
	return len(flow.Messages) == 0
}

func (flow *messagesFlow) LastMessage() Message {
	if len(flow.Messages) == 0 {
		return nil
	}

	var lastMessage Message

	for _, message := range flow.Messages {
		if message.SentAt().After(lastMessage.SentAt()) {
			lastMessage = message
		}
	}

	return lastMessage
}
