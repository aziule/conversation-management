package conversation

import (
	"errors"
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
	ErrCannotUnmarshalBson         = errors.New("Can't unmarshal BSON")
)

// Repository is the main interface for accessing conversation-related objects
type Repository interface {
	FindLatestConversation(user *User) (*Conversation, error)
	SaveConversation(conversation *Conversation) error
	FindUserByFbId(fbId string) (*User, error)
	InsertUser(user *User) error
}

// MessageWithType is the struct grouping a message along with its type.
// Its purpose is to be used with mongodb so that we can unmarshal any Message interface easily.
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

// StartConversation initialises a new conversation
func StartConversation() *Conversation {
	return &Conversation{
		Status:    StatusOngoing,
		Messages:  nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddMessage is called when a new message needs to be added to the conversation
func (conversation *Conversation) AddMessage(message *UserMessage) {
	conversation.Messages = append(
		conversation.Messages,
		&MessageWithType{message.Type(), message},
	)
}

// IsNew tells us if the conversation is a new one
func (conversation *Conversation) IsNew() bool {
	return len(conversation.Messages) == 0
}

// SetBSON converts the MessageWithType's message, of type interface, to the corresponding Message struct.
// We have moved the SetBSON method here for now, as it was faster to develop. In the future,
// we may think about extracting it completely to the "conversation/mongo" package.
func (m *MessageWithType) SetBSON(raw bson.Raw) error {
	decodedType := struct {
		Type MessageType `bson:"type"`
	}{}

	err := raw.Unmarshal(&decodedType)

	if err != nil {
		log.Infof("Could not unmarshal BSON: %s", err)
		return ErrCannotUnmarshalBson
	}

	switch decodedType.Type {
	case MessageFromUser:
		decodedMessage := struct {
			Message *UserMessage `bson:"message"`
		}{}
		raw.Unmarshal(&decodedMessage)
		m.Message = decodedMessage.Message
		break
	default:
		log.WithField("type", decodedType.Type).Infof("Could not unmarshal BSON: unhandled message type")
		return ErrCannotUnmarshalBson
	}

	m.Type = decodedType.Type

	return nil
}
