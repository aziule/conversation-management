package conversation

import "github.com/aziule/conversation-management/core/nlp"

// StepHandler handles any step given
type StepHandler interface {
	// CanValidate will check if the step is valid using the given ParsedData
	CanValidate(step *Step, data nlp.ParsedData) bool

	// Process will process the step and take the relevant actions
	Process(step *Step, data nlp.ParsedData) error
}

// Story is the main structure for user stories, which are basically
// a flow of steps to validate and process.
type Story struct {
	Name          string
	StartingSteps []*Step
}

// Step is the main structure for the various steps taken within a single story.
// Each step consists of a name and a set of expectations, in terms
// of intents or entities (data).
// Each step links to the next ones, until there are no more steps,
// in which case we can consider the Story as done.
type Step struct {
	Name             string
	ExpectedIntent   nlp.Intent
	ExpectedEntities []nlp.Entity
	NextSteps        []*Step
}

// FindStep returns the step with the provided name, if found.
// Returns nil if no step is found.
func (s *Story) FindStep(name string) *Step {
	for _, step := range s.StartingSteps {
		if step.Name == name {
			return step
		}

		found := step.findSubStep(name)

		if found != nil {
			return found
		}
	}

	return nil
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
