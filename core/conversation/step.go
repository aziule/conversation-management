package conversation

import (
	"github.com/aziule/conversation-management/core/nlu"
)

// Step is a structure acting as a link between other steps in a given Story
type Step struct {
	name             string
	next             []*Step // Next available steps
	expectedIntent   *nlu.Intent
	expectedEntities []*nlu.Entity
}

// NewStep is the constructor method
func NewStep(name string, next []*Step) *Step {
	return &Step{
		name: name,
		next: next,
	}
}

// Getters
func (step *Step) Name() string                    { return step.name }
func (step *Step) Next() []*Step                   { return step.next }
func (step *Step) ExpectedIntent() *nlu.Intent     { return step.expectedIntent }
func (step *Step) ExpectedEntities() []*nlu.Entity { return step.expectedEntities }

// IsExpectingIntent will tell us if the step is expecting an intent in order to be valid
func (step *Step) IsExpectingIntent() bool {
	return step.expectedIntent != nil
}

// IsExpectingEntity will tell us if the step is expecting one or more entities in order to be valid
func (step *Step) IsExpectingEntity() bool {
	return len(step.expectedEntities) > 0
}
