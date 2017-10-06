package conversation

type Context struct {
	currentStep *Step
}

// Retrieve the current context for a user
func getContext(user User) (*Context, error) {
	return nil, nil
}

// Getters
func (context *Context) CurrentStep() *Step { return context.currentStep }
