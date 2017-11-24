package conversation

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrStoryRepositoryNotFound = errors.New("StoryRepository not found")

	// storyRepositoryBuilders stores the available StoryRepository builders
	storyRepositoryBuilders = make(map[string]StoryRepositoryBuilder)
)

// StoryRepositoryBuilder is the interface describing a builder for StoryRepository
type StoryRepositoryBuilder func() StoryRepository

// RegisterStoryRepositoryBuilder adds a new StoryRepositoryBuilder to the list of available builders
func RegisterStoryRepositoryBuilder(name string, builder StoryRepositoryBuilder) {
	_, registered := storyRepositoryBuilders[name]

	if registered {
		log.WithField("name", name).Warning("StoryRepositoryBuilder already registered, ignoring")
	}

	storyRepositoryBuilders[name] = builder
}

// NewStoryRepository tries to create a StoryRepository using the available builders.
// Returns ErrStoryRepositoryNotFound if the repository builder isn't found.
func NewStoryRepository(name string) (StoryRepository, error) {
	storyRepositoryBuilder, ok := storyRepositoryBuilders[name]

	if !ok {
		return nil, ErrStoryRepositoryNotFound
	}

	return storyRepositoryBuilder(), nil
}

// StoryRepository is the repository responsible for fetching our stories
type StoryRepository interface {
	FindAll() ([]*Story, error)
}

// Story is the main structure for user stories, which are basically
// a flow of steps to step in and process.
type Story struct {
	Name          string
	StartingSteps []*Step
}

// NewStory is our constructor method for Story
func NewStory(name string, startingSteps []*Step) *Story {
	return &Story{
		Name:          name,
		StartingSteps: startingSteps,
	}
}

// AddStartingStep attaches a new starting step to the story
func (s *Story) AddStartingStep(step *Step) {
	s.StartingSteps = append(s.StartingSteps, step)
}

// FindStep returns the step with the provided name, if found.
// Returns nil if no step is found.
// @todo: return an error instead and handle the not found with an ErrNotFound
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
