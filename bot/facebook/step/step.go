package step

import "errors"

type StepHandlerFunc func()

// stepHandlers contains all of the available StepHandlerFunc
var stepHandlers = map[string]StepHandlerFunc{}

// StepHandler represents an interface that will be used to handle a specific step
type StepHandler struct {}

// is where we register our static Facebook step handlers for now
func init() {
	stepHandlers["step1"] = Step1
	stepHandlers["step2"] = Step2
}

// HandleStep is the method we call when trying to handle a step
func HandleStep(name string) error {
	stepHandlerFunc, available := stepHandlers[name]

	if ! available {
		// @todo: better handling of errors
		return errors.New("The step could not be handled")
	}

	stepHandlerFunc()

	return nil
}