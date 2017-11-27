package wit

import "github.com/aziule/conversation-management/core/nlp"

// witApi is the struct used to access data from WIT
type witApi struct {
}

func newWitApi() nlp.Repository {
	return &witApi{}
}

// GetIntents is the method returning all of the available intents.
// It is required in order to implement the nlp.Repository interface.
func (api *witApi) GetIntents() ([]*nlp.Intent, error) {
	return []*nlp.Intent{
		nlp.NewIntent("test"),
		nlp.NewIntent("test2"),
	}, nil
}

func init() {
	nlp.RegisterRepositoryBuilder("wit", newWitApi)
}
