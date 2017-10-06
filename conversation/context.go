package conversation

type Context struct {
	currentStep *Step
}

func GetContext(user *User) (*Context, error) {

}

// Getters
func (context *Context) CurrentStep() *Step { return context.currentStep }
