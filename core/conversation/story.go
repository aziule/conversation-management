package conversation

// Story is the main structure for user stories, which are basically
// a flow of steps to validate and process.
type Story struct {
	Name          string
	StartingSteps []*Step
}

// StoryRepository is the repository responsible for fetching our stories
type StoryRepository interface {
	FindAll() ([]*Story, error)
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
