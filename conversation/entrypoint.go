package conversation

// Represents the main entrypoint for a chatbot, with all of the available stories
type Entrypoint struct {
	stories []*Story
}

// Constructor method
func NewEntryPoint(stories []*Story) *Entrypoint {
	return &Entrypoint{
		stories: stories,
	}
}

// Getters
func (entrypoint *Entrypoint) Stories() []*Story { return entrypoint.stories }