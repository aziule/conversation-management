package conversation

import (
	"errors"

	"github.com/aziule/conversation-management/core/nlp"
	log "github.com/sirupsen/logrus"
)

// Step is the main structure for the various steps taken within a single story.
// Each step consists of a name and a set of expectations, in terms
// of intents or entities (data).
// Each step links to the next ones, until there are no more steps,
// in which case we can consider the Story as done.
type Step struct {
	Name             string
	ExpectedIntent   string
	ExpectedEntities []string
	NextSteps        []*Step
}

// NewStep is our constructor method for Step
func NewStep(name string, expectedIntent string, expectedEntities []string, nextSteps []*Step) *Step {
	return &Step{
		Name:             name,
		ExpectedIntent:   expectedIntent,
		ExpectedEntities: expectedEntities,
		NextSteps:        nextSteps,
	}
}

// AddNextStep adds a new step to the list of the next available steps
func (s *Step) AddNextStep(step *Step) {
	s.NextSteps = append(s.NextSteps, step)
}

// IsLastStep tells us if a step is the last one.
//
// Simply put, if a step does not have next steps, then it's
// considered as being the last one.
func (s *Step) IsLastStep() bool {
	return len(s.NextSteps) == 0
}

// findSubStep looks for sub steps within a step, and returns the one matching the provided name.
// Returns nil if no sub step is found.
func (s *Step) findSubStep(name string) *Step {
	if len(s.NextSteps) == 0 {
		return nil
	}

	for _, subStep := range s.NextSteps {
		if subStep.Name == name {
			return subStep
		}

		return subStep.findSubStep(name)
	}

	return nil
}

// StepProcessFunc is a func responsible for handling a given step
type StepProcessFunc func(step *Step, data *nlp.ParsedData) error

// StepsProcessMap maps steps names to their process func
type StepsProcessMap map[string]StepProcessFunc

// StepHandler is the struct responsible for handling steps for a bot
type StepHandler struct {
	processMap StepsProcessMap
}

// NewStepHandler is the constructor method for StepHandler
func NewStepHandler(processMap StepsProcessMap) *StepHandler {
	return &StepHandler{
		processMap: processMap,
	}
}

// CanStepIn tries to see if the NLP data meets the step's requirements
// in order to process the step. It will check if the expected intent / entities
// are present in the NLP data, and return true or false accordingly.
//
// However, this method does not check the data itself. It only checks
// for its presence, not its validity.
// @todo: needs testing
func (h *StepHandler) CanStepIn(step *Step, data *nlp.ParsedData) bool {
	// Case 1: NLP data provides an intent but it's not the same name
	if data.Intent != nil && step.ExpectedIntent != data.Intent.Name {
		return false
	}

	// Case 2: NLP data does not provide an intent but we are expecting one
	if data.Intent == nil && step.ExpectedIntent != "" {
		return false
	}

	if len(step.ExpectedEntities) > 0 {
		for _, expectedEntity := range step.ExpectedEntities {
			hasEntity := false

			for _, providedEntity := range data.Entities {
				if providedEntity.Entity.Name() == expectedEntity {
					hasEntity = true
					log.Debugf("Has entity: %s", expectedEntity)
				}
			}

			if !hasEntity {
				log.Debugf("Missing entity to step in: %s", expectedEntity)
				return false
			}
		}
	}

	return true
}

// Process will process the step using its associated StepProcessFunc.
// Returns an error if there is no associated StepProcessFunc or
// for any other processing reason.
func (h *StepHandler) Process(step *Step, data *nlp.ParsedData) error {
	fn, ok := h.processMap[step.Name]

	if !ok {
		// @todo: handle this case and log
		return errors.New("Cannot handle")
	}

	return fn(step, data)
}
