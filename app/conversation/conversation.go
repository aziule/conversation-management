package conversation

import (
	"errors"
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

// @todo: manage the internal state (current step, status, etc.)
// Conversation is the struct that will handle our conversations between
// the bot and the various users.
type Conversation struct {
	Id       bson.ObjectId `bson:"_id"`
	Status   Status        `bson:"status"`
	Messages []Message     `bson:"messages"`
}

// NewConversation is the constructor for Conversation.
// The initial status is ongoing
func NewConversation() *Conversation {
	return &Conversation{
		Status:   StatusOngoing,
		Messages: nil,
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
