package conversation

type Context struct {
	currentStep *Step
}

func getContext(user User) (*Context, error) {
	return nil, nil
}

// Getters
func (context *Context) CurrentStep() *Step { return context.currentStep }
