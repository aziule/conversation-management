package conversation

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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

type MessageWithType struct {
	Type    MessageType `bson:"type"`
	Message Message     `bson:"message"`
}

// @todo: manage the internal state (current step, status, etc.)
// Conversation is the struct that will handle our conversations between
// the bot and the various users.
type Conversation struct {
	Id        bson.ObjectId      `bson:"_id"`
	Status    Status             `bson:"status"`
	Messages  []*MessageWithType `bson:"messages"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
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
	conversation.Messages = append(
		conversation.Messages,
		&MessageWithType{message.Type(), message},
	)
}

// IsNew tells us if the conversation is a new one
func (conversation *Conversation) IsNew() bool {
	return len(conversation.Messages) == 0
}

// SetBSON converts the MessageWithType's message to a correct Message struct
// as the MessageWithType.Message is an interface.
func (ml *MessageWithType) SetBSON(raw bson.Raw) error {
	fmt.Println("Here in SetBSON")

	decodedType := struct {
		Type MessageType `bson:"type"`
	}{}

	raw.Unmarshal(&decodedType)

	switch decodedType.Type {
	case MessageFromUser:
		decodedMessage := struct {
			Message *UserMessage `bson:"message"`
		}{}
		raw.Unmarshal(&decodedMessage)
		ml.Message = decodedMessage.Message
		break
	default:
		// @todo: handle and log
		log.WithField("type", decodedType.Type).Error("Unhandled message type to decode")
		return errors.New("Unhandled case")
	}

	ml.Type = decodedType.Type

	return nil
}
