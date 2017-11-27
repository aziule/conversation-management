package nlp

import (
	"errors"

	"github.com/aziule/conversation-management/core/utils"
	log "github.com/sirupsen/logrus"
)

var (
	ErrApiNotFound = errors.New("Api not found")

	// apiBuilders stores the available Api builders
	apiBuilders = make(map[string]ApiBuilder)
)

// ApiBuilder is the interface describing a builder for Api
type ApiBuilder func(conf utils.BuilderConf) (Api, error)

// RegisterApiBuilder adds a new ApiBuilder to the list of available builders
func RegisterApiBuilder(name string, builder ApiBuilder) {
	_, registered := apiBuilders[name]

	if registered {
		log.WithField("name", name).Warning("ApiBuilder already registered, ignoring")
	}

	apiBuilders[name] = builder
}

// NewApi tries to create a Api using the available builders.
// Returns ErrApiNotFound if the api builder isn't found.
func NewApi(name string, conf utils.BuilderConf) (Api, error) {
	apiBuilder, ok := apiBuilders[name]

	if !ok {
		return nil, ErrApiNotFound
	}

	api, err := apiBuilder(conf)

	if err != nil {
		return nil, err
	}

	return api, nil
}

// Api is the main interface used to get / store NLP data
type Api interface {
	GetIntents() ([]*Intent, error)
}
