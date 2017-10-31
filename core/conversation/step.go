package conversation

type StepLoader interface {
	LoadSteps() *Step
}

type StepHandler interface {
	Handle(step *Step) error
}

type Step struct {
	Name string
}
