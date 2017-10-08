package conversation

// Steps are a link between other steps in a Story
type Step struct {
	name string
	next []*Step // Next available steps
}

// Constructor method
func NewStep(name string, next []*Step) *Step {
	return &Step{
		name: name,
		next: next,
	}
}

// Getters
func (step *Step) Name() string { return step.name }
func (step *Step) Next() []*Step { return step.next }
