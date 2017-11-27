package wit

import (
	"net/http"

	"github.com/aziule/conversation-management/core/nlp"
	"github.com/aziule/conversation-management/core/utils"
)

// witApi is the struct used to make calls to WIT
type witApi struct {
	client *http.Client
}

// newWitApi creates a new witApi using the given conf
func newWitApi(conf utils.BuilderConf) (nlp.Api, error) {
	client, ok := utils.GetParam(conf, "client").(*http.Client)

	if !ok {
		return nil, utils.ErrInvalidOrMissingParam("client")
	}

	return &witApi{
		client: client,
	}, nil
}

// GetIntents gets the list of intents from WIT
func (api *witApi) GetIntents() ([]*nlp.Intent, error) {
	return []*nlp.Intent{
		nlp.NewIntent("test"),
		nlp.NewIntent("test2"),
	}, nil
}

func init() {
	nlp.RegisterApiBuilder("wit", newWitApi)
}
