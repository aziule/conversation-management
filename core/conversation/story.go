package conversation

// Story is the entry point to a given task
type Story struct {
	name          string
	startingSteps []*Step
}

// NewStory is the constructor method
func NewStory(name string, startingSteps []*Step) *Story {
	return &Story{
		name:          name,
		startingSteps: startingSteps,
	}
}

// Getters
func (story *Story) Name() string           { return story.name }
func (story *Story) StartingSteps() []*Step { return story.startingSteps }
