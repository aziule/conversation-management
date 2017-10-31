package wit

import "github.com/aziule/conversation-management/core/conversation"

type witStepLoader struct {
}

func NewWitStepLoader() conversation.StepLoader {
	return &witStepLoader{}
}

func (loader *witStepLoader) LoadSteps() error {

}
