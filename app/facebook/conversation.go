package facebook

import (
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// conversationHandler is the struct responsible for handling Facebook conversations
type conversationHandler struct {
	conversationRepository conversation.Repository
}

// NewConversationHandler is the constructor method for conversationHandler
func NewConversationHandler(conversationRepository conversation.Repository) *conversationHandler {
	return &conversationHandler{
		conversationRepository: conversationRepository,
	}
}

// GetConversation tries to return a Facebook conversation between a given user and the bot.
// If there is an ongoing conversation, then it will return it.
// If this is the first conversation or the previous one is marked as done, then it will create a new one.
func (h *conversationHandler) GetConversation(user *conversation.User) (*conversation.Conversation, error) {
	// => this will help with consolidated users (fb + slack + anything)
	c, err := h.conversationRepository.FindLatestConversation(user)

	if err != nil {
		if err != conversation.ErrNotFound {
			return nil, err
		}

		log.WithField("user", user).Info("Starting a first conversation")

		// The conversation was not found: start a new one
		c = conversation.StartConversation()
	}

	// Start a new conversation if the previous one is over
	if c.Status == conversation.StatusOver {
		log.WithField("user", user).Info("Starting a new conversation")

		c = conversation.StartConversation()
	}

	return c, nil
}

func (h *conversationHandler) HandleStep(step *conversation.Step) error {
	return nil
}

// GetUser tries to find an existing user using the id provided as the facebook id.
// If it does not find any user then it will create a new one using the facebook id.
func (h *conversationHandler) GetUser(id string) (*conversation.User, error) {
	user, err := h.conversationRepository.FindUserByFbId(id)

	if err != nil && err != conversation.ErrNotFound {
		return nil, err
	}

	if user == nil {
		log.WithField("fbId", id).Infof("Inserting a new user")

		user = &conversation.User{
			Id:   bson.NewObjectId(),
			FbId: id,
		}

		// Insert the user
		err = h.conversationRepository.InsertUser(user)

		if err != nil {
			// @todo: handle this case and return something to the user
			log.WithField("fbId", id).Info("Could not insert the user")
			return nil, err
		}
	}

	return user, nil
}

// stepHandler is the struct responsible for handling steps for a Facebook bot
type stepHandler struct {
}

func (h *stepHandler) CanValidate(step *conversation.Step, data nlp.ParsedData) bool {
	return false
}

func (h *stepHandler) Process(step *conversation.Step, data nlp.ParsedData) error {
	return nil
}
