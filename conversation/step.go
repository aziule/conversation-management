package conversation

// Steps are a link between other steps in a Story
type Step struct {
	name string
	next []*Step // Next available steps
}

// Getters
func (step *Step) Name() string { return step.name }
func (step *Step) Next() string { return step.next }
