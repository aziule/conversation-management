package wit

import (
	"net/http"

	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/core/utils"
)

const baseUrl = "http://api.wit.ai"

// witApi is the struct used to make calls to Wit
type witApi struct {
	client *http.Client
}

// newWitApi creates a new witApi using the given conf
func newWitApi(conf utils.BuilderConf) (interface{}, error) {
	client, ok := utils.GetParam(conf, "client").(*http.Client)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("client")
	}

	return &witApi{
		client: client,
	}, nil
}

// GetEntities gets the list of entities from Wit
func (api *witApi) GetEntities() ([]nlp.Entity, error) {
	return []nlp.Entity{
		nlp.NewIntEntity("test", 0, 123),
	}, nil
}

func init() {
	nlp.RegisterApiBuilder("wit", newWitApi)
}
