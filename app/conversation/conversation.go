package conversation

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Status string

var (
	StatusOngoing           Status = "ongoing"
	StatusHumanIntervention Status = "human"
	StatusOver              Status = "over"
	ErrNotFound                    = errors.New("Not found")
)

// Repository is the main interface for accessing conversation-related objects
type Repository interface {
	FindLatestConversation(user *User) (*Conversation, error)
	SaveConversation(conversation *Conversation) error
	FindUser(userId string) (*User, error)
	InsertUser(user *User) error
}

// MessagesList contains the list of messages for a conversation.
// It contains every messages sent from and to the bot.
type MessagesList []Message

// @todo: manage the internal state (current step, status, etc.)
// Conversation is the struct that will handle our conversations between
// the bot and the various users.
type Conversation struct {
	Id        bson.ObjectId `bson:"_id"`
	Status    Status        `bson:"status"`
	Messages  MessagesList  `bson:"messages"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

// NewConversation is the constructor for Conversation.
// The initial status is ongoing
func NewConversation() *Conversation {
	return &Conversation{
		Status:    StatusOngoing,
		Messages:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Received is called when a new message is received
func (conversation *Conversation) Received(message *UserMessage) {
	conversation.Messages = append(conversation.Messages, message)
}

// IsNew tells us if the conversation is a new one
func (conversation *Conversation) IsNew() bool {
	return len(conversation.Messages) == 0
}
