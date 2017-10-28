package conversation

import (
	"time"

	"github.com/aziule/conversation-management/app/nlp"
)

type MessageType string

const (
	MessageFromUser MessageType = "from-user"
	MessageFromBot  MessageType = "from-bot"
)

type Message interface {
	Text() string
	SentAt() time.Time
	Type() MessageType
}

// message is the struct that contains the base information about a generic message
type message struct {
	Text   string      `bson:"text"`
	SentAt time.Time   `bson:"sent_at"`
	Type   MessageType `bson:"type"`
}

// newMessage is the private constructor method for message
func newMessage(text string, sentAt time.Time, messageType MessageType) *message {
	return &message{
		Text:   text,
		SentAt: sentAt,
		Type:   messageType,
	}
}

// UserMessage represents a message received from a user
type UserMessage struct {
	*message
	Sender     *User
	ParsedData *nlp.ParsedData `bson:"parsed_data"`
}

// NewUserMessage is the constructor method for UserMessage
func NewUserMessage(text string, sentAt time.Time, sender *User, parsedData *nlp.ParsedData) *UserMessage {
	return &UserMessage{
		newMessage(text, sentAt, MessageFromUser),
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

func (msg *UserMessage) Type() MessageType {
	return msg.message.Type
}

type BotMessage struct {
	message
	RepliesTo *UserMessage
}
