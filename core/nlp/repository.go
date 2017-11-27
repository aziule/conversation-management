package nlp

import (
	"github.com/aziule/conversation-management/core/utils"
)

const nlpRepositoryBuilderPrefix = "nlp_repository_"

// RegisterRepositoryBuilder registers a new service builder using a package-level prefix using a package-level prefix
func RegisterRepositoryBuilder(name string, builder utils.ServiceBuilder) {
	utils.RegisterServiceBuilder(nlpRepositoryBuilderPrefix+name, builder)
}

// NewRepository tries to create a Repository using the available builders.
// Returns ErrRepositoryNotFound if the repository builder isn't found.
func NewRepository(name string, conf utils.BuilderConf) (Repository, error) {
	repositoryBuilder, err := utils.GetServiceBuilder(nlpRepositoryBuilderPrefix + name)

	if err != nil {
		return nil, err
	}

	repository, err := repositoryBuilder(conf)

	if err != nil {
		return nil, err
	}

	return repository.(Repository), nil
}

// Repository is the main interface used to get / store NLP data
type Repository interface {
	GetIntents() ([]*Intent, error)
}
