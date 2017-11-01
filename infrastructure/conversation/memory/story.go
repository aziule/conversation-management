package memory

import "github.com/aziule/conversation-management/core/conversation"

type inMemoryStepRepository struct {
}

func NewStepRepository() *inMemoryStepRepository {
	return &inMemoryStepRepository{}
}

func (r *inMemoryStepRepository) FindAll() ([]*conversation.Story, error) {
	return nil, nil
}
