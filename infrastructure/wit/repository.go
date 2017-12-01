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
	entities, err := repository.api.GetEntities()

	if err != nil {
		// @todo: log
		return nil, err
	}

	intents := []*nlp.Intent{}

	for _, entity := range entities {
		if entity.Type == nlp.IntentEntity {
			intents = append(intents, nlp.NewIntent(entity.Name))
		}
	}

	return intents, nil
}

func init() {
	nlp.RegisterRepositoryBuilder("wit", newWitRepository)
}
