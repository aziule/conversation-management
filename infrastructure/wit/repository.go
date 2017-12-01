package wit

import (
	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/core/utils"
)

// witRepository is the struct used to access data from Wit
type witRepository struct {
	api nlp.Api
}

// newWitRepository instantiates a new witRepository using the given conf
func newWitRepository(conf utils.BuilderConf) (interface{}, error) {
	api, ok := utils.GetParam(conf, "api").(nlp.Api)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("api")
	}

	return &witRepository{
		api: api,
	}, nil
}

// GetIntents is the method returning all of the available intents.
// It is required in order to implement the nlp.Repository interface.
func (repository *witRepository) GetIntents() ([]*nlp.Intent, error) {
	return repository.api.GetIntents()
}

// GetEntities is the method returning all of the available entities.
// It is required in order to implement the nlp.Repository interface.
func (repository *witRepository) GetEntities() ([]*nlp.Entity, error) {
	return repository.api.GetEntities()
}

func init() {
	nlp.RegisterRepositoryBuilder("wit", newWitRepository)
}
