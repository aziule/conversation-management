package conversation

// Story are the entry points to a given task
type Story struct {
	name string
	startingSteps []*Step
}

// Getters
func (story *Story) Name() string { return story.name }
func (story *Story) StartingSteps() []*Step { return story.startingSteps }
