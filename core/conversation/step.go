package conversation

import "github.com/aziule/conversation-management/core/nlp"

// StepHandler handles any step given
type StepHandler interface {
	// CanStepIn will check if the NLP data is enough to allow the step to be processed
	CanStepIn(step *Step, data *nlp.ParsedData) bool

	// Process will process the step and take the relevant actions
	Process(step *Step, data *nlp.ParsedData) error
}

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
