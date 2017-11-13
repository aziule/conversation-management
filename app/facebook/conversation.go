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
	stepHandler            *conversation.StepHandler
	conversationRepository conversation.Repository
	storyRepository        conversation.StoryRepository
}

// newConversationHandler is the constructor method for conversationHandler
func newConversationHandler(pm conversation.StepsProcessMap, cr conversation.Repository, sr conversation.StoryRepository) *conversationHandler {
	return &conversationHandler{
		stepHandler:            conversation.NewStepHandler(pm),
		conversationRepository: cr,
		storyRepository:        sr,
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
		// @todo: handle: save the user message here?
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
		log.WithField("story", story).Debugf("Trying to step in story")
		if startingStep != nil {
			break
		}

		for _, step := range story.StartingSteps {
			if h.stepHandler.CanStepIn(step, data) {
				log.WithField("step", step).Debugf("Stepping in")

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

		return errors.New("Handle this. Don't forget to save the conversation with the message")
	}

	// Process the step
	err = h.stepHandler.Process(startingStep, data)

	if err != nil {
		// @todo: handle this, and save the conversation's message otherwise it's lost
		return errors.New("nope")
	}

	// Update the conversation
	return nil
}

// tryProgressInStory is the method being called when a conversation is ongoing and we try to progress
// within the current story.
func (h *conversationHandler) tryProgressInStory(data *nlp.ParsedData, c *conversation.Conversation) error {
	stories, err := h.storyRepository.FindAll()

	if err != nil {
		// @todo: log, and save conversation
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

	if currentStep == nil {
		// @todo: return a correct error message.
		// @todo: handle this case and see how we can prevent
		// a conversation from being blocked.
		return errors.New("Could not find any step")
	}

	canStepIn := h.stepHandler.CanStepIn(currentStep, data)

	log.Info("Stepping in: %s", canStepIn)

	return nil
}
