package facebook

import (
	"errors"
	"net/http"

	"github.com/aziule/conversation-management/core/conversation"
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/infrastructure/facebook/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// conversationHandler is the struct responsible for handling Facebook conversations
type conversationHandler struct {
	stepHandler            *conversation.StepHandler
	conversationRepository conversation.Repository
	storyRepository        conversation.StoryRepository
	nlpParser              nlp.Parser
	fbApi                  api.FacebookApi
}

// newConversationHandler is the constructor method for conversationHandler
func newConversationHandler(pm conversation.StepsProcessMap, cr conversation.Repository, sr conversation.StoryRepository, p nlp.Parser, a api.FacebookApi) *conversationHandler {
	return &conversationHandler{
		stepHandler:            conversation.NewStepHandler(pm),
		conversationRepository: cr,
		storyRepository:        sr,
		nlpParser:              p,
		fbApi:                  a,
	}
}

// MessageReceived is the implementation of ConversationHandler.MessageReceived method.
// It handles the whole conversation processing logic for Facebook bots.
//
// - Parsing the request
// - Validating the message
// - Parsing NLP
// - Managing the conversation flow
// - Answering the user
// - Modifying the conversation's status
func (h *conversationHandler) MessageReceived(r *http.Request) error {
	receivedMessage, err := h.fbApi.ParseRequestMessageReceived(r)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("message", receivedMessage).Errorf("Could not parse the received message: %s", err)
		return err
	}

	user, err := h.getUser(receivedMessage.SenderId)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("user", receivedMessage.SenderId).Errorf("Could not find the user: %s", err)
		return err
	}

	c, err := h.getConversation(user)

	if err != nil {
		// @todo: handle this case and return something to the user
		log.WithField("user", user).Infof("Could not get the conversation: %s", err)
		return err
	}

	log.WithField("conversation", c).Debug("Conversation fetched")

	userMessage := conversation.NewUserMessage(
		receivedMessage.Text,
		receivedMessage.SentAt,
		user,
		nil,
	)

	c.AddMessage(userMessage)

	h.conversationRepository.SaveConversation(c)

	if receivedMessage.Nlp == nil {
		// @todo: handle this case: parse the text using the NLP parser
		log.Errorf("No data to parse")
		return err
	}

	parsedData, err := h.nlpParser.ParseNlpData(receivedMessage.Nlp)

	if err != nil {
		// @todo: handle this case and return something to the user. Make sure the
		// conversation is saved with the message. For example, we could think
		// about adding a flag to the message, like:
		// - could_not_parse_nlp
		// - could_not_process
		// - something_else
		// - ...
		// => gives more context and allows us to save data & understand it even
		// though errors occur.
		// @todo: save the conversation
		log.WithField("nlp", receivedMessage.Nlp).Errorf("Could not parse NLP data: %s", err)
		return err
	}

	userMessage.ParsedData = parsedData

	log.WithField("data", parsedData).Debug("Data parsed from message")

	h.conversationRepository.SaveConversation(c)

	err = h.processData(parsedData, c)

	if err != nil {
		log.WithFields(log.Fields{
			"data":         parsedData,
			"conversation": c,
		}).Errorf("Could not process the data: %s", err)
	}

	return nil
}

// processData is the method responsible for taking actions on a conversation using the provided NLP data
func (h *conversationHandler) processData(data *nlp.ParsedData, c *conversation.Conversation) error {
	var err error

	if c.CurrentStep == "" {
		log.WithField("c", c).Info("Try starting a new story")
		err = h.tryStartStory(data, c)
	} else {
		log.WithField("c", c).Info("Try progressing in the current story")
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

	return h.processStep(c, startingStep, data)
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
		log.WithFields(log.Fields{
			"data":         data,
			"conversation": c,
		}).Error("The conversation's current step does not exist in the stories")

		// @todo: return a correct error message.
		// @todo: handle this case and see how we can prevent
		// a conversation from being blocked.
		return errors.New("Could not find any step")
	}

	var nextStep *conversation.Step

	for _, step := range currentStep.NextSteps {
		if h.stepHandler.CanStepIn(step, data) {
			log.WithField("step", step).Debugf("Stepping in")

			nextStep = step
			break
		}
	}

	if nextStep == nil {
		log.WithFields(log.Fields{
			"data":         data,
			"conversation": c,
		}).Info("Cannot progress in story")

		return errors.New("Handle this. Don't forget to save the conversation with the message")
	}

	return h.processStep(c, nextStep, data)
}

// processStep processes a single step, according to the fact that we should
// be able, at that stage, to step in the step.
//
// So make sure to call step.CanStepIn and that the result is true
// before calling this method.
func (h *conversationHandler) processStep(c *conversation.Conversation, s *conversation.Step, data *nlp.ParsedData) error {
	// Process the step
	log.WithFields(log.Fields{
		"step": s,
		"data": data,
	}).Info("Processing step")

	err := h.stepHandler.Process(s, data)

	if err != nil {
		log.Errorf("Could not process the step: %s", err)
		// @todo: handle this, and save the conversation's message otherwise it's lost
		return errors.New("nope")
	}

	// Update the conversation's state
	c.CurrentStep = s.Name
	h.conversationRepository.SaveConversation(c)

	return nil
}

// getConversation tries to return a Facebook conversation between a given user and the bot.
// If there is an ongoing conversation, then it will return it.
// If this is the first conversation or the previous one is marked as done, then it will create a new one.
func (h *conversationHandler) getConversation(user *conversation.User) (*conversation.Conversation, error) {
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

// getUser tries to find an existing user using the id provided as the facebook id.
// If it does not find any user then it will create a new one using the facebook id.
func (h *conversationHandler) getUser(id string) (*conversation.User, error) {
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
