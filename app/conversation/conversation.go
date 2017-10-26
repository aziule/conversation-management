package conversation

import "errors"

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

// Conversation is the struct that will handle our conversations between
// the bot and the various users.
type Conversation struct {
	Status       Status
	messagesFlow *messagesFlow
}

// NewConversation is the constructor for Conversation.
// The initial status is ongoing
func NewConversation() *Conversation {
	return &Conversation{
		Status: StatusOngoing,
	}
}

func (conversation *Conversation) Received(message *UserMessage) {
	conversation.messagesFlow.Messages = append(conversation.messagesFlow.Messages, message)
}
