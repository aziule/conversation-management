package nlp

import (
	"github.com/aziule/conversation-management/core/utils"
)

const nlpApiBuilderPrefix = "nlp_api_"

// RegisterRepositoryBuilder registers a new service builder using a package-level prefix
func RegisterApiBuilder(name string, builder utils.ServiceBuilder) {
	utils.RegisterServiceBuilder(nlpApiBuilderPrefix+name, builder)
}

// NewApi tries to create a Api using the available builders.
// Returns ErrApiNotFound if the api builder isn't found.
func NewApi(name string, conf utils.BuilderConf) (Api, error) {
	apiBuilder, err := utils.GetServiceBuilder(nlpApiBuilderPrefix + name)

	if err != nil {
		return nil, err
	}

	api, err := apiBuilder(conf)

	if err != nil {
		return nil, err
	}

	return api.(Api), nil
}

// Api is the main interface used to get / store NLP data
type Api interface {
	GetEntities() ([]Entity, error)
}
