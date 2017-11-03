package facebook

import (
	"errors"
	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// conversationHandler is the struct responsible for handling Facebook conversations
type conversationHandler struct {
	stepHandler            conversation.StepHandler
	conversationRepository conversation.Repository
	storyRepository        conversation.StoryRepository
}

// newConversationHandler is the constructor method for conversationHandler
func newConversationHandler(conversationRepository conversation.Repository, storyRepository conversation.StoryRepository) *conversationHandler {
	return &conversationHandler{
		stepHandler:            newStepHandler(),
		conversationRepository: conversationRepository,
		storyRepository:        storyRepository,
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
		c = conversation.CreateNewConversation()
	}

	// Start a new conversation if the previous one is over
	if c.Status == conversation.StatusOver {
		log.WithField("user", user).Info("Starting a new conversation")

		c = conversation.CreateNewConversation()
	}

	return c, nil
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

// ProcessData is the method responsible for taking actions on a conversation using the provided NLP data
func (h *conversationHandler) ProcessData(data *nlp.ParsedData, c *conversation.Conversation) error {
	var err error

	if c.CurrentStep == "" {
		log.Debug("Try starting a new story")
		err = h.tryStartStory(data, c)
	} else {
		log.Debug("Try progressing in the current story")
		err = h.tryProgressInStory(data, c)
	}

	if err != nil {
		// @todo: handle
		return err
	}

	return nil
}

// tryStartStory will try to start a new story using the provided NLP data.
// It will go through the available stories and see if any step can be initiated.
func (h *conversationHandler) tryStartStory(data *nlp.ParsedData, c *conversation.Conversation) error {
	stories, err := h.storyRepository.FindAll()

	if err != nil {
		// @todo: log
		return errors.New("Cannot load stories")
	}

	var startingStep *conversation.Step

	for _, story := range stories {
		log.WithField("story", story).Debugf("Trying to validate story")
		if startingStep != nil {
			break
		}

		for _, step := range story.StartingSteps {
			log.WithField("step", step).Debugf("Trying to validate starting step")
			if h.stepHandler.CanValidate(step, data) {
				startingStep = step
				break
			}
		}
	}

	if startingStep == nil {
		log.WithFields(log.Fields{
			"data":         data,
			"conversation": c,
		}).Info("Cannot start a story")

		return errors.New("Handle this")
	}

	// Update the conversation
	return nil
}

func (h *conversationHandler) tryProgressInStory(data *nlp.ParsedData, c *conversation.Conversation) error {
	stories, err := h.storyRepository.FindAll()

	if err != nil {
		// @todo: log
		return errors.New("Cannot load stories")
	}

	var currentStep *conversation.Step

	// Find the current step of the conversation
	for _, story := range stories {
		step := story.FindStep(c.CurrentStep)

		if step != nil {
			currentStep = step
			break
		}
	}

	canValidate := h.stepHandler.CanValidate(currentStep, data)

	log.Info("Validating conversation: %s", canValidate)

	return nil
}

// stepHandler is the struct responsible for handling steps for a Facebook bot
// @todo: add mapping: step's name <=> handler func
type stepHandler struct {
}

// newStepHandler is the constructor method for stepHandler
func newStepHandler() *stepHandler {
	return &stepHandler{}
}

// CanValidate tries to validate a step given the NLP data.
// It will check for the expected intent / entities and return true or false accordingly.
func (h *stepHandler) CanValidate(step *conversation.Step, data *nlp.ParsedData) bool {

	return false
}

func (h *stepHandler) Process(step *conversation.Step, data *nlp.ParsedData) error {
	return nil
}
